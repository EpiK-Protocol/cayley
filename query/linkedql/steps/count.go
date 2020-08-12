package steps

import (
	"github.com/cayleygraph/quad/voc"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/epik-protocol/epik-gateway-backend/query/linkedql"
	"github.com/epik-protocol/epik-gateway-backend/query/path"
)

func init() {
	linkedql.Register(&Count{})
}

var _ linkedql.PathStep = (*Count)(nil)

// Count corresponds to .count().
type Count struct {
	From linkedql.PathStep `json:"from"`
}

// Description implements Step.
func (s *Count) Description() string {
	return "resolves to the number of the resolved values of the from step"
}

// BuildPath implements linkedql.PathStep.
func (s *Count) BuildPath(qs graph.QuadStore, ns *voc.Namespaces) (*path.Path, error) {
	fromPath, err := s.From.BuildPath(qs, ns)
	if err != nil {
		return nil, err
	}
	return fromPath.Count(), nil
}
