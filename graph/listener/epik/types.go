package epik

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

var EmptyTSK = TipSetKey{}

// The length of a block header CID in bytes.
var blockHeaderCIDLen int

func init() {
	c, err := cid.V1Builder{Codec: cid.DagCBOR, MhType: multihash.BLAKE2B_MIN + 31}.Sum([]byte{})
	if err != nil {
		panic(err)
	}
	blockHeaderCIDLen = len(c.Bytes())
}

type TipSet struct {
	Cids   []cid.Cid
	Height abi.ChainEpoch
}

func (ts *TipSet) Key() TipSetKey {
	if ts == nil {
		return EmptyTSK
	}
	return NewTipSetKey(ts.Cids...)
}

type TipSetKey struct {
	// The empty key has value "".
	value string
}

func NewTipSetKey(cids ...cid.Cid) TipSetKey {
	encoded := encodeKey(cids)
	return TipSetKey{string(encoded)}
}

func (k TipSetKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Cids())
}

func (k *TipSetKey) UnmarshalJSON(b []byte) error {
	var cids []cid.Cid
	if err := json.Unmarshal(b, &cids); err != nil {
		return err
	}
	k.value = string(encodeKey(cids))
	return nil
}

func encodeKey(cids []cid.Cid) []byte {
	buffer := new(bytes.Buffer)
	for _, c := range cids {
		// bytes.Buffer.Write() err is documented to be always nil.
		_, _ = buffer.Write(c.Bytes())
	}
	return buffer.Bytes()
}

// Cids returns a slice of the CIDs comprising this key.
func (k TipSetKey) Cids() []cid.Cid {
	cids, err := decodeKey([]byte(k.value))
	if err != nil {
		panic("invalid tipset key: " + err.Error())
	}
	return cids
}

func decodeKey(encoded []byte) ([]cid.Cid, error) {
	// To avoid reallocation of the underlying array, estimate the number of CIDs to be extracted
	// by dividing the encoded length by the expected CID length.
	estimatedCount := len(encoded) / blockHeaderCIDLen
	cids := make([]cid.Cid, 0, estimatedCount)
	nextIdx := 0
	for nextIdx < len(encoded) {
		nr, c, err := cid.CidFromBytes(encoded[nextIdx:])
		if err != nil {
			return nil, err
		}
		cids = append(cids, c)
		nextIdx += nr
	}
	return cids, nil
}

type EpikClient interface {
	ChainHead(context.Context) (*TipSet, error)
	ChainGetTipSetByHeight(context.Context, abi.ChainEpoch, TipSetKey) (*TipSet, error)
	ChainGetBlockMessages(ctx context.Context, blockCid cid.Cid) (*BlockMessages, error)
	ClientFindData(ctx context.Context, root cid.Cid) error
}

type EpikClientStruct struct {
	Internal struct {
		ChainHead              func(context.Context) (*TipSet, error)                              `perm:"read"`
		ChainGetTipSetByHeight func(context.Context, abi.ChainEpoch, TipSetKey) (*TipSet, error)   `perm:"read"`
		ChainGetBlockMessages  func(ctx context.Context, blockCid cid.Cid) (*BlockMessages, error) `perm:"read"`
		ClientFindData         func(ctx context.Context, root cid.Cid) error                       `perm:"read"`
	}
}

func (e *EpikClientStruct) ChainHead(ctx context.Context) (*TipSet, error) {
	return e.Internal.ChainHead(ctx)
}

func (e *EpikClientStruct) ChainGetTipSetByHeight(ctx context.Context, epoch abi.ChainEpoch, tsk TipSetKey) (*TipSet, error) {
	return e.Internal.ChainGetTipSetByHeight(ctx, epoch, tsk)
}

func (e *EpikClientStruct) ChainGetBlockMessages(ctx context.Context, blockCid cid.Cid) (*BlockMessages, error) {
	return e.Internal.ChainGetBlockMessages(ctx, blockCid)
}

func (e *EpikClientStruct) ClientFindData(ctx context.Context, root cid.Cid) error {
	return e.Internal.ClientFindData(ctx, root)
}
