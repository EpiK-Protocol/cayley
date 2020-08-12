package httpgraph

import (
	"net/http"

	"github.com/epik-protocol/epik-gateway-backend/graph"
)

type QuadStore interface {
	graph.QuadStore
	ForRequest(r *http.Request) (graph.QuadStore, error)
}
