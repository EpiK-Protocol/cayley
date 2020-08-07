package epik

import (
	"bytes"
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

const (
	FileDownloading = 0
	FileDownloaded  = 1
)

type FileData struct {
	Status int
	Url    string
	Cid    cid.Cid
}
