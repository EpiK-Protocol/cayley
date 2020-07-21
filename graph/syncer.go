package graph

import (
	"context"
	"fmt"
)

var (
	ErrSyncerNotRegistred = fmt.Errorf("this syncer is not registered")
)

var syncerRegistry = make(map[string]SyncerRegistration)

type HandleDeltaFunc func([]Delta) error

type SyncerRegistration struct {
	NewSyncerFunc func(QuadStore) (Syncer, error)
}

type Syncer interface {
	// at - ChainEpoch
	RetrieveAll(ctx context.Context, at int64, handler HandleDeltaFunc)
	// since, to - ChainEpoch
	RetrieveDelta(ctx context.Context, since, to int64, handler HandleDeltaFunc)
	Start()
	Stop()
}

func RegisterSyncer(name string, register SyncerRegistration) {
	if register.NewSyncerFunc == nil {
		panic("NewSyncerFunc must not be nil")
	}

	if _, found := syncerRegistry[name]; found {
		panic(fmt.Sprintf("Already registered Syncer %q.", name))
	}
	syncerRegistry[name] = register
}

func NewSyncer(name string, store QuadStore) (Syncer, error) {
	r, registered := syncerRegistry[name]
	if !registered {
		return nil, ErrSyncerNotRegistred
	}
	return r.NewSyncerFunc(store)
}
