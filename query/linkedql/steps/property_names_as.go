package steps

import (
	"github.com/cayleygraph/quad/voc"
	"github.com/epik-protocol/gateway/graph"
	"github.com/epik-protocol/gateway/query/linkedql"
	"github.com/epik-protocol/gateway/query/path"
)

func init() {
	linkedql.Register(&PropertyNamesAs{})
}

var _ linkedql.PathStep = (*PropertyNamesAs)(nil)

// PropertyNamesAs corresponds to .propertyNamesAs().
type PropertyNamesAs struct {
	From linkedql.PathStep `json:"from"`
	Tag  string            `json:"tag"`
}

// Description implements Step.
func (s *PropertyNamesAs) Description() string {
	return "tags the list of predicates that are pointing out from a node."
}

// BuildPath implements linkedql.PathStep.
func (s *PropertyNamesAs) BuildPath(qs graph.QuadStore, ns *voc.Namespaces) (*path.Path, error) {
	fromPath, err := s.From.BuildPath(qs, ns)
	if err != nil {
		return nil, err
	}
	return fromPath.SavePredicates(false, s.Tag), nil
}
