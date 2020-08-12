package steps

import (
	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/voc"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/epik-protocol/epik-gateway-backend/graph/iterator"
	"github.com/epik-protocol/epik-gateway-backend/query/linkedql"
	"github.com/epik-protocol/epik-gateway-backend/query/path"
)

func init() {
	linkedql.Register(&LessThanEquals{})
}

var _ linkedql.PathStep = (*LessThanEquals)(nil)

// LessThanEquals corresponds to lte().
type LessThanEquals struct {
	From  linkedql.PathStep `json:"from"`
	Value quad.Value        `json:"value"`
}

// Description implements Step.
func (s *LessThanEquals) Description() string {
	return "Less than equals filters out values that are not less than or equal given value"
}

// BuildPath implements linkedql.PathStep.
func (s *LessThanEquals) BuildPath(qs graph.QuadStore, ns *voc.Namespaces) (*path.Path, error) {
	fromPath, err := s.From.BuildPath(qs, ns)
	if err != nil {
		return nil, err
	}
	return fromPath.Filter(iterator.CompareLTE, linkedql.AbsoluteValue(s.Value, ns)), nil
}
