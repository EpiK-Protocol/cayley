package epik

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/epik-protocol/epik-gateway-backend/clog"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
	multiaddr "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/spf13/viper"
)

type RepoType int

const (
	flagEpikAddr  = "datasource.options.address"
	flagEpikToken = "datasource.options.token"

	envFullNodeInfo = "FULLNODE_API_INFO"
)

type APIInfo struct {
	Addr  multiaddr.Multiaddr
	Token []byte
}

type EpikClient interface {
	ChainHead(context.Context) (*TipSet, error)
	ChainGetTipSetByHeight(context.Context, abi.ChainEpoch, TipSetKey) (*TipSet, error)
	ChainGetBlockMessages(ctx context.Context, blockCid cid.Cid) (*BlockMessages, error)
	ClientQuery(ctx context.Context, roots []cid.Cid) ([]FileResp, error)
	ChainGetParentMessages(ctx context.Context, blockCid cid.Cid) ([]APIMessage, error)
	ChainGetMessage(context.Context, cid.Cid) (*Message, error)
}

type EpikClientStruct struct {
	Internal struct {
		ChainHead              func(context.Context) (*TipSet, error)                              `perm:"read"`
		ChainGetTipSetByHeight func(context.Context, abi.ChainEpoch, TipSetKey) (*TipSet, error)   `perm:"read"`
		ChainGetBlockMessages  func(ctx context.Context, blockCid cid.Cid) (*BlockMessages, error) `perm:"read"`
		ClientQuery            func(ctx context.Context, roots []cid.Cid) ([]FileResp, error)      `perm:"read"`
		ChainGetParentMessages func(ctx context.Context, blockCid cid.Cid) ([]APIMessage, error)   `perm:"read"`
		ChainGetMessage        func(context.Context, cid.Cid) (*Message, error)                    `perm:"read"`
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

func (e *EpikClientStruct) ChainGetParentMessages(ctx context.Context, blockCid cid.Cid) ([]APIMessage, error) {
	return e.Internal.ChainGetParentMessages(ctx, blockCid)
}

func (e *EpikClientStruct) ChainGetMessage(ctx context.Context, blockCid cid.Cid) (*Message, error) {
	return e.Internal.ChainGetMessage(ctx, blockCid)
}

func (e *EpikClientStruct) ClientQuery(ctx context.Context, roots []cid.Cid) ([]FileResp, error) {
	return e.Internal.ClientQuery(ctx, roots)
}

func (a APIInfo) DialArgs() (string, error) {
	_, addr, err := manet.DialArgs(a.Addr)
	return "ws://" + addr + "/rpc/v0", err
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", string(a.Token))
		return headers
	}
	clog.Warningf("API Token not set and requested, capabilities might be limited.")
	return nil
}

func NewEpikClient() (EpikClient, jsonrpc.ClientCloser, error) {

	ainfo, err := GetAPIInfo()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get API info: %w", err)
	}

	addr, err := ainfo.DialArgs()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get DialArgs: %w", err)
	}

	var res EpikClientStruct
	closer, err := jsonrpc.NewMergeClient(addr, "EpiK",
		[]interface{}{
			&res.Internal,
		},
		ainfo.AuthHeader(),
	)

	return &res, closer, err
}

func GetAPIInfo() (APIInfo, error) {

	info := APIInfo{}

	if env, ok := os.LookupEnv(envFullNodeInfo); ok {
		sp := strings.SplitN(env, ":", 2)
		if len(sp) != 2 {
			clog.Warningf("invalid env(%s) value, missing token or address", envFullNodeInfo)
		} else {
			ma, err := multiaddr.NewMultiaddr(sp[1])
			if err != nil {
				return APIInfo{}, fmt.Errorf("could not parse multiaddr from env(%s): %w", envFullNodeInfo, err)
			}
			info.Addr = ma
			info.Token = []byte(sp[0])
		}
	}

	flagAddr := viper.GetString(flagEpikAddr)
	flagToken := viper.GetString(flagEpikToken)
	if len(flagAddr) > 0 {
		ma, err := multiaddr.NewMultiaddr(flagAddr)
		if err != nil {
			return APIInfo{}, fmt.Errorf("could not parse multiaddr(%s) from flag: %w", flagAddr, err)
		}
		info.Addr = ma
	}
	if len(flagToken) > 0 {
		info.Token = []byte(flagToken)
	}
	return info, nil
}
