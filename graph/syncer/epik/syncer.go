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
	"github.com/ipfs/go-cid"
	"github.com/spf13/viper"
)

const (
	SyncerType = "epik"

	keyEpikAddr = "epik.address"

	syncDuration = 30 * time.Second

	readBatchSize = 10000
)

func init() {
	graph.RegisterSyncer(SyncerType, graph.SyncerRegistration{
		NewSyncerFunc: func(store graph.QuadStore) (graph.Syncer, error) {
			return newSyncer(store)
		},
	})
}

type Syncer struct {
	quit  chan struct{}
	start sync.Once
	wg    sync.WaitGroup

	client EpikClient
	store  graph.QuadStore
}

func newSyncer(store graph.QuadStore) (*Syncer, error) {

	return &Syncer{
		quit:  make(chan struct{}),
		store: store,
	}, nil
}

func (s *Syncer) Start() {
	s.start.Do(func() {
		s.wg.Add(1)
		go s.startSyncer()
	})
}

func (s *Syncer) startSyncer() {
	var (
		err    error
		close  func()
		ticker = time.NewTicker(syncDuration)
	)
	s.client, close, err = NewEpikClient(viper.GetString(keyEpikAddr), makeRequestHeader())
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

			if err = s.doSync(local.Epoch+1, remote); err != nil {
				clog.Errorf("failed to sync epick from %d to %d, error is: %v", local.Epoch+1, remote, err)
			}
		}
	}
}

func (s *Syncer) Stop() {
	close(s.quit)
	s.wg.Wait()
}

// TODO: does chain selection matters?
func (s *Syncer) doSync(current, end int64) error {
	for ; current <= end; current++ {
		select {
		case <-s.quit:
			break
		default:
		}

		adds, deletes, err := s.client.GetChange(context.TODO(), current)
		if err != nil {
			clog.Errorf("failed to get change at epoch %d, error is: %v", current, err)
			return err
		}

		// process delete
		procDels, err := parseDeletes(current, deletes)
		if err != nil {
			clog.Errorf("failed to parse deleted cids at epoch %d, error is: %v", current, err)
			return err
		}
		if err = s.store.ApplyDeltas(procDels, graph.IgnoreOpts{IgnoreMissing: true}); err != nil {
			clog.Errorf("failed to apply deletes at epoch %d, error is: %v", current, err)
			return err
		}

		// process add
		m, err := s.client.GetObjects(context.TODO(), adds)
		if err != nil {
			clog.Errorf("failed to get objects at epoch %d, error is: %v", current, err)
			return err
		}

		procAdds, err := parseAdds(current, m)
		if err != nil {
			clog.Errorf("failed to parse added cids at epoch %d, error is: %v", current, err)
			return err
		}

		if err = s.store.ApplyDeltas(procAdds, graph.IgnoreOpts{IgnoreDup: true}); err != nil {
			clog.Errorf("failed to apply adds at epoch %d, error is: %v", current, err)
			return err
		}
	}
	return nil
}

func parseDeletes(epoch int64, ids []cid.Cid) ([]graph.Delta, error) {
	deltas := make([]graph.Delta, 0, len(ids))
	for _, id := range ids {
		deltas = append(deltas, graph.Delta{
			Cid:    id.String(),
			Action: graph.Delete,
			Epoch:  epoch,
		})
	}
	return deltas, nil
}

func parseAdds(epoch int64, m map[cid.Cid][]byte) ([]graph.Delta, error) {
	deltas := make([]graph.Delta, 0, len(m))
	for id, data := range m {
		// See load.go:83
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
				Cid:    id.String(),
				Quad:   q,
				Action: graph.Add,
				Epoch:  epoch,
			})
		}
		qr.Close()
	}
	return deltas, nil
}
