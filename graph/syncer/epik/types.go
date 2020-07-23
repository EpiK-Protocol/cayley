package epik

import (
	"context"

	"github.com/ipfs/go-cid"
)

type EpikClient interface {
	GetBestEpoch(ctx context.Context) (int64, error)
	GetChange(ctx context.Context, epoch int64) (adds []cid.Cid, deletes []cid.Cid, err error)
	GetObjects(ctx context.Context, ids []cid.Cid) (map[cid.Cid][]byte, error)
}

type EpikClientStruct struct {
	Internal struct {
		GetBestEpoch func(ctx context.Context) (int64, error)
		GetChange    func(ctx context.Context, epoch int64) (adds []cid.Cid, deletes []cid.Cid, err error)
		GetObjects   func(ctx context.Context, ids []cid.Cid) (map[cid.Cid][]byte, error)
	}
}

func (e *EpikClientStruct) GetBestEpoch(ctx context.Context) (int64, error) {
	return e.Internal.GetBestEpoch(ctx)
}

func (e *EpikClientStruct) GetChange(ctx context.Context, epoch int64) (adds []cid.Cid, deletes []cid.Cid, err error) {
	return e.Internal.GetChange(ctx, epoch)
}

func (e *EpikClientStruct) GetObjects(ctx context.Context, ids []cid.Cid) (map[cid.Cid][]byte, error) {
	return e.Internal.GetObjects(ctx, ids)
}
