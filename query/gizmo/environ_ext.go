package gizmo

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/epik-protocol/epik-gateway-backend/query/path"
)

// S is a shorthand for Search.
func (g *graphObject) NewS(call goja.FunctionCall) goja.Value {
	return g.NewSearch(call)
}

// Search uses a Third-Party search engine(e.g. Elasticsearch) for fast node search.
// If no search engine or arguments provided, it will fall back to Vertex.
func (g *graphObject) NewSearch(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 || g.s.gns == nil {
		return g.NewVertex(call)
	}
	keys := make(map[string]struct{}, len(call.Arguments))
	for _, o := range exportArgs(call.Arguments) {
		key, ok := o.(string)
		if !ok {
			return throwErr(g.s.vm, fmt.Errorf("not a string for Search: %T", o))
		}
		keys[key] = struct{}{}
	}
	qv, err := g.s.gns.SearchNodes(keys)
	if err != nil {
		return throwErr(g.s.vm, err)
	}
	return g.s.vm.ToValue(&pathObject{
		s:      g.s,
		finals: true,
		path:   path.StartMorphism(qv...),
	})
}
