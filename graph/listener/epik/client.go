package epik

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/epik-protocol/gateway/clog"
	"github.com/filecoin-project/go-jsonrpc"
	multiaddr "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/spf13/viper"
)

type RepoType int

const (
	flagEpikAddr  = "epik.address"
	flagEpikToken = "epik.token"

	envFullNodeInfo = "FULLNODE_API_INFO"
)

type APIInfo struct {
	Addr  multiaddr.Multiaddr
	Token []byte
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
