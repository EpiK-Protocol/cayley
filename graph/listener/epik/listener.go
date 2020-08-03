package epik

import (
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	"github.com/cayleygraph/quad/nquads"
	"github.com/epik-protocol/gateway/clog"
	"github.com/epik-protocol/gateway/graph"
)

const (
	ListenerType = "epik"

	syncDuration = 30 * time.Second

	readBatchSize = 10000
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
	var (
		err    error
		close  func()
		ticker = time.NewTicker(syncDuration)
	)
	s.client, close, err = NewEpikClient()
	if err != nil {
		clog.Fatalf("failed to init epik client: %v", err)
	}

	defer func() {
		err := recover()
		if err != nil {
			clog.Errorf("panic: %v", err)
		}
		close()
		ticker.Stop()
	}()

	for {
		select {
		case <-s.quit:
			clog.Infof("epik syncer stopped")
			return
		case <-ticker.C:
			local, err := s.store.Stats(context.TODO(), false)
			if err != nil {
				clog.Errorf("failed to get quadstore stats: %v", err)
				continue
			}

			remote, err := s.client.GetBestEpoch(context.TODO())
			if err != nil {
				clog.Errorf("failed to get best epoch: %v", err)
				continue
			}

			if err = s.syncDeltas(local.Epoch+1, remote); err != nil {
				clog.Errorf("failed to sync epick from %d to %d, error is: %v", local.Epoch+1, remote, err)
			}
		}
	}
}

func (s *Listener) Stop() {
	close(s.quit)
	s.wg.Wait()
}

func (s *Listener) syncDeltas(current, end int64) error {
	for ; current <= end; current++ {
		select {
		case <-s.quit:
			break
		default:
		}

		msgs, err := s.client.GetMessages(context.TODO(), current)
		if err != nil {
			clog.Errorf("failed to get change at epoch %d, error is: %v", current, err)
			return err
		}

		// delete
		deltas, err := parseDeletes(msgs.Deletes)
		if err != nil {
			clog.Errorf("failed to parse deleted cids at epoch %d, error is: %v", current, err)
			return err
		}

		// add
		objs, err := s.client.GetObjects(context.TODO(), msgs.Adds)
		if err != nil {
			clog.Errorf("failed to get objects at epoch %d, error is: %v", current, err)
			return err
		}

		adds, err := parseAdds(objs)
		if err != nil {
			clog.Errorf("failed to parse added cids at epoch %d, error is: %v", current, err)
			return err
		}

		deltas = append(deltas, adds...)
		if err = s.store.ApplyDeltas(current, deltas, graph.IgnoreOpts{IgnoreDup: true, IgnoreMissing: true}); err != nil {
			clog.Errorf("failed to apply adds at epoch %d, error is: %v", current, err)
			return err
		}
	}
	return nil
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
