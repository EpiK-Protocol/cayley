package epik

import (
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
)

func makeRequestHeader() http.Header {
	return make(http.Header)
}

func NewEpikClient(addr string, reqHeader http.Header) (EpikClient, jsonrpc.ClientCloser, error) {
	var res EpikClientStruct
	closer, err := jsonrpc.NewMergeClient(addr, "Epik",
		[]interface{}{
			&res.Internal,
		},
		reqHeader,
	)
	return &res, closer, err
}
