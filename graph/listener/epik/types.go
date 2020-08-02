package epik

import (
	"context"
)

type Messages struct {
	Adds    []string // cid
	Deletes []string // cid
}

type EpikClient interface {
	GetBestEpoch(ctx context.Context) (int64, error)
	GetMessages(ctx context.Context, epoch int64) (msg Messages, err error)
	GetObjects(ctx context.Context, cids []string) (map[string][]byte, error)
}

type EpikClientStruct struct {
	Internal struct {
		GetBestEpoch func(ctx context.Context) (int64, error)                            `perm:"read"`
		GetMessages  func(ctx context.Context, epoch int64) (msg Messages, err error)    `perm:"read"`
		GetObjects   func(ctx context.Context, cids []string) (map[string][]byte, error) `perm:"read"`
	}
}

func (e *EpikClientStruct) GetBestEpoch(ctx context.Context) (int64, error) {
	return e.Internal.GetBestEpoch(ctx)
}

func (e *EpikClientStruct) GetMessages(ctx context.Context, epoch int64) (msg Messages, err error) {
	return e.Internal.GetMessages(ctx, epoch)
}

func (e *EpikClientStruct) GetObjects(ctx context.Context, cids []string) (map[string][]byte, error) {
	if len(cids) == 0 {
		return make(map[string][]byte), nil
	}
	return e.Internal.GetObjects(ctx, cids)
}
