package epik

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cayleygraph/quad/nquads"
	"github.com/epik-protocol/gateway/clog"
	"github.com/epik-protocol/gateway/graph"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
)

const (
	ListenerType = "epik"

	syncDuration = 10 * time.Minute
)

func init() {
	graph.RegisterListener(ListenerType, graph.ListenerRegistration{
		NewListenerFunc: func(store graph.QuadStore) (graph.Listener, error) {
			return newListener(store)
		},
	})
}

type Listener struct {
	quit  chan struct{}
	start sync.Once
	wg    sync.WaitGroup

	client EpikClient
	store  graph.QuadStore
}

func newListener(store graph.QuadStore) (*Listener, error) {

	return &Listener{
		quit:  make(chan struct{}),
		store: store,
	}, nil
}

func (s *Listener) Start() {
	s.start.Do(func() {
		s.wg.Add(1)
		go s.listen()
	})
}

func (s *Listener) listen() {

	defer func() {
		err := recover()
		if err != nil {
			clog.Errorf("panic: %v", err)
		}
		s.wg.Done()
	}()

	ticker := time.NewTicker(syncDuration)
	defer ticker.Stop()

	var (
		err   error
		close func()
	)
	s.client, close, err = NewEpikClient()
	if err != nil {
		clog.Fatalf("failed to init epik client: %v", err)
	}
	defer close()

	ctx, cancelFunc := context.WithCancel(context.Background())

	for {
		select {
		case <-s.quit:
			cancelFunc()
			clog.Infof("epik syncer stopped")
			return
		case <-ticker.C:
			local, err := s.store.Stats(ctx, false)
			if err != nil {
				clog.Errorf("failed to get quadstore stats: %v", err)
				continue
			}

			head, err := s.client.ChainHead(ctx)
			if err != nil {
				clog.Errorf("failed to get head epoch: %v", err)
				continue
			}
			remote := int64(head.Height)
			if remote <= 1 {
				continue
			}

			if err = s.syncDeltas(ctx, local.Epoch+1, remote-1); err != nil {
				clog.Errorf("failed to sync deltas from %d to %d, error is: %v", local.Epoch+1, remote-1, err)
			}
		}
	}
}

func (s *Listener) Stop() {
	close(s.quit)
	s.wg.Wait()
}

func (s *Listener) syncDeltas(ctx context.Context, start, end int64) error {
	tss, err := s.getTipSets(ctx, start, end)
	if err != nil {
		return err
	}

	for _, ts := range tss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// get tipset messages
		msgs, err := s.getTipSetMessages(ctx, ts)
		if err != nil {
			return err
		}

		// filter storage deals
		cids := make([]cid.Cid, 0, len(msgs))
		for _, msg := range msgs {
			if msg.To != builtin.StorageMarketActorAddr ||
				msg.Method != builtin.MethodsMarket.PublishStorageDeals {
				continue
			}

			var params market.PublishStorageDealsParams
			if err := params.UnmarshalCBOR(bytes.NewReader(msg.Params)); err != nil {
				clog.Errorf("failed to unmarshal params at tipset %d, error is: %v", ts.Height, err)
				return err
			}
			cids = append(cids, params.RootCID)
		}

		// add to QuadStore
		// TODO: delete
		// deltas, err := parseDeletes(msgs.Deletes)
		// if err != nil {
		// 	clog.Errorf("failed to parse deleted cids at epoch %d, error is: %v", current, err)
		// 	return err
		// }

		quadAdds, err := s.retrieveFiles(ctx, cids)
		if err != nil {
			clog.Errorf("failed to wait at epoch %d, error is: %v", ts.Height, err)
			return err
		}

		adds, err := parseAdds(quadAdds)
		if err != nil {
			clog.Errorf("failed to parse added cids at epoch %d, error is: %v", ts.Height, err)
			return err
		}

		if err = s.store.ApplyDeltas(int64(ts.Height), adds, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: true}); err != nil {
			clog.Errorf("failed to apply adds at epoch %d, error is: %v", ts.Height, err)
			return err
		}
	}
	// just set epoch to "end"
	return s.store.ApplyDeltas(end, nil, graph.IgnoreOpts{})
}

func (s *Listener) retrieveFiles(ctx context.Context, cids []cid.Cid) (map[string][]byte, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	downloader := newDownloader()
loop:
	for {
		var unfinished []cid.Cid
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			if len(cids) == 0 {
				break loop
			}
			fds, err := s.client.ClientQuery(ctx, cids)
			if err != nil {
				return nil, err
			}

			for _, fd := range fds {
				switch fd.Status {
				case FileDownloaded:
					if err := downloader.download(ctx, fd.Root.String(), fd.Url); err != nil {
						return nil, err
					}
				case FileDownloading:
					unfinished = append(unfinished, fd.Root)
				default:
					return nil, fmt.Errorf("unexpected file status: %d", fd.Status)
				}
			}
		}
		cids = unfinished
	}
	downloader.wait()
	if downloader.failed > 0 {
		return nil, fmt.Errorf("file downloader failed %d times", downloader.failed)
	}
	return downloader.result, nil
}

func (s *Listener) getTipSets(ctx context.Context, start, end int64) ([]*TipSet, error) {
	from := EmptyTSK
	tsl := list.New()

	for start <= end {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		ts, err := s.client.ChainGetTipSetByHeight(ctx, abi.ChainEpoch(end), from)
		if err != nil {
			clog.Errorf("failed to get tipset at epoch %d, error is: %v", end, err)
			return nil, err
		}
		if int64(ts.Height) < start {
			break
		}
		tsl.PushFront(ts)
		from = ts.Key()
		end = int64(ts.Height) - 1
	}
	r := make([]*TipSet, 0, tsl.Len())
	for e := tsl.Front(); e != nil; e.Next() {
		r = append(r, e.Value.(*TipSet))
	}
	return r, nil
}

func (s *Listener) getTipSetMessages(ctx context.Context, ts *TipSet) ([]*Message, error) {
	// get tipset messages
	msgs := make([]*Message, 0, 100)
	for _, bcid := range ts.Cids {
		bm, err := s.client.ChainGetBlockMessages(ctx, bcid)
		if err != nil {
			clog.Errorf("failed to get block messages at tipset %d, error is: %v", ts.Height, err)
			return nil, err
		}
		for _, m := range bm.BlsMessages {
			msgs = append(msgs, m)
		}
		for _, m := range bm.SecpkMessages {
			msgs = append(msgs, &m.Message)
		}
	}
	return msgs, nil
}

func parseDeletes(cids []string) ([]graph.Delta, error) {
	deltas := make([]graph.Delta, 0, len(cids))
	for _, cid := range cids {
		if len(cid) == 0 {
			return nil, graph.ErrInvalidCid
		}
		deltas = append(deltas, graph.Delta{
			Cid:    cid,
			Action: graph.Delete,
		})
	}
	return deltas, nil
}

// map value maybe empty for expiration
func parseAdds(m map[string][]byte) ([]graph.Delta, error) {
	deltas := make([]graph.Delta, 0, len(m))
	for cid, data := range m {
		if len(cid) == 0 {
			return nil, graph.ErrInvalidCid
		}
		// load.go:83
		qr := nquads.NewReader(bytes.NewReader(data), false)
		for {
			q, err := qr.ReadQuad()
			if err != nil {
				if err == io.EOF {
					break
				}
				qr.Close()
				return nil, err
			}
			deltas = append(deltas, graph.Delta{
				Cid:    cid,
				Quad:   q,
				Action: graph.Add,
			})
		}
		qr.Close()
	}
	return deltas, nil
}
