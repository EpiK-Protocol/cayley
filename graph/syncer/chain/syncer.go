package chainsyncer

import (
	"context"
	"sync"

	"github.com/epik-protocol/gateway/clog"
	"github.com/epik-protocol/gateway/graph"
)

const SyncerType = "chain"

func init() {
	graph.RegisterSyncer(SyncerType, graph.SyncerRegistration{
		NewSyncerFunc: func(store graph.QuadStore) (graph.Syncer, error) {
			return newSyncer(store)
		},
	})
}

// TODO: add chain retrieve interface
type Syncer struct {
	quit  chan struct{}
	start sync.Once
	wg    sync.WaitGroup

	store graph.QuadStore
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
	defer func() {
		err := recover()
		if err != nil {
			clog.Errorf("failed to start chain syncer: %v", err)
		}
	}()
}

func (s *Syncer) Stop() {
	close(s.quit)
	s.wg.Wait()
}

func (s *Syncer) RetrieveAll(ctx context.Context, at int64, handler graph.HandleDeltaFunc) {

}

func (s *Syncer) RetrieveDelta(ctx context.Context, since, to int64, handler graph.HandleDeltaFunc) {

}
