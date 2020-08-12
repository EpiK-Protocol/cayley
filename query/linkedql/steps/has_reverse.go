package steps

import (
	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/voc"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/epik-protocol/epik-gateway-backend/query/linkedql"
	"github.com/epik-protocol/epik-gateway-backend/query/path"
)

func init() {
	linkedql.Register(&HasReverse{})
}

var _ linkedql.PathStep = (*HasReverse)(nil)

// HasReverse corresponds to .hasR().
type HasReverse struct {
	From     linkedql.PathStep      `json:"from"`
	Property *linkedql.PropertyPath `json:"property"`
	Values   []quad.Value           `json:"values"`
}

// Description implements Step.
func (s *HasReverse) Description() string {
	return "is the same as Has, but sets constraint in reverse direction."
}

// BuildPath implements linkedql.PathStep.
func (s *HasReverse) BuildPath(qs graph.QuadStore, ns *voc.Namespaces) (*path.Path, error) {
	fromPath, err := s.From.BuildPath(qs, ns)
	if err != nil {
		return nil, err
	}
	viaPath, err := s.Property.BuildPath(qs, ns)
	if err != nil {
		return nil, err
	}
	return fromPath.HasReverse(viaPath, linkedql.AbsoluteValues(s.Values, ns)...), nil
}
