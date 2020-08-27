package graph

import (
	"errors"
	"strings"

	"github.com/cayleygraph/quad"
)

type NewSearcherFunc func(QuadStore, Options) (Searcher, error)

var searcherRegistry = make(map[string]NewSearcherFunc)

type Searcher interface {
	SearchNodes(keys map[string]struct{}) ([]quad.Value, error)
	IndexNodes([]Node) error
	SetMeta(map[string]interface{}) error
	GetMeta(key string) (interface{}, error)
	Close() error
}

func RegisterSearcher(name string, newFunc NewSearcherFunc) {
	if _, found := searcherRegistry[name]; found {
		panic("already registered Searcher " + name)
	}
	searcherRegistry[name] = newFunc
}

func NewSearcher(name string, qs QuadStore, opts Options) (Searcher, error) {
	newFunc, hasNew := searcherRegistry[name]
	if !hasNew {
		return nil, errors.New("node searcher: name '" + name + "' is not registered")
	}
	return newFunc(qs, opts)
}

type Node struct {
	// internal ID
	ID      int64  `json:"id"`
	Content string `json:"content"`
	Type    string `json:"type"` // "string", "iri", "bnode"
}

// TODO:
func (n *Node) NotIndex() bool {
	if strings.HasPrefix(n.Content, "data:image/") {
		return true
	}
	return false
}
