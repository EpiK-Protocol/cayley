package epik

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	big2 "github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/ipfs/go-cid"
)

type BigInt = big2.Int

type APIMessage struct {
	Cid     cid.Cid
	Message *Message
}

type Message struct {
	Version int64

	To   address.Address
	From address.Address

	Nonce uint64

	Value BigInt

	GasPrice BigInt
	GasLimit int64

	Method abi.MethodNum
	Params []byte
}

type SignedMessage struct {
	Message   Message
	Signature crypto.Signature
}

type BlockMessages struct {
	BlsMessages   []*Message
	SecpkMessages []*SignedMessage

	Cids []cid.Cid
}
