package steps

import (
	"github.com/cayleygraph/quad/voc"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/epik-protocol/epik-gateway-backend/query/linkedql"
	"github.com/epik-protocol/epik-gateway-backend/query/path"
)

func init() {
	linkedql.Register(&ReversePropertyNames{})
}

var _ linkedql.PathStep = (*ReversePropertyNames)(nil)

// ReversePropertyNames corresponds to .reversePropertyNames().
type ReversePropertyNames struct {
	From linkedql.PathStep `json:"from"`
}

// Description implements Step.
func (s *ReversePropertyNames) Description() string {
	return "gets the list of predicates that are pointing in to a node."
}

// BuildPath implements linkedql.PathStep.
func (s *ReversePropertyNames) BuildPath(qs graph.QuadStore, ns *voc.Namespaces) (*path.Path, error) {
	fromPath, err := s.From.BuildPath(qs, ns)
	if err != nil {
		return nil, err
	}
	return fromPath.InPredicates(), nil
}
