package es

import (
	"fmt"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	"github.com/stretchr/testify/assert"
)

func TestSetMeta(t *testing.T) {
	es, err := NewSearcher(nil, graph.Options{
		"addresses": []string{"http://localhost:9200"},
	})
	assert.Nil(t, err)

	err = es.SetMeta(map[string]interface{}{
		"maxid": 100,
	})
	assert.Nil(t, err)

	value, err := es.GetMeta("maxid")
	assert.Nil(t, err)
	t.Logf("after setting maxid: %v", value)
}

func TestIndexNodes(t *testing.T) {
	es, err := NewSearcher(nil, graph.Options{
		"addresses": []string{"http://localhost:9200"},
	})
	assert.Nil(t, err)

	nodes := []graph.Node{
		{
			ID:      0,
			Content: "alice buy mass",
			Type:    "string",
		},
		{
			ID:      1,
			Content: "bob sail 2 mass",
			Type:    "string",
		},
		{
			ID:      2,
			Content: "<signal mass>",
			Type:    "iri",
		},
		{
			ID:      3,
			Content: "_:nothing bob",
			Type:    "bnode",
		},
	}
	err = es.IndexNodes(nodes)
	assert.Nil(t, err)
}

func TestSearchNodes(t *testing.T) {
	es, err := NewSearcher(nil, graph.Options{
		"addresses": []string{"http://localhost:9200"},
	})
	assert.Nil(t, err)

	qds, err := es.SearchNodes(map[string]struct{}{
		"buy":  {},
		"Sail": {},
	})
	assert.Nil(t, err)

	for _, qd := range qds {
		fmt.Println(qd.Native())
	}
}

func TestCheckIndices(t *testing.T) {

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	assert.Nil(t, err)

	err = CheckAndCreateIndex(es, IndexNode, index_node_mapping)
	assert.Nil(t, err)
	err = CheckAndCreateIndex(es, IndexMeta, index_meta_mapping)
	assert.Nil(t, err)
}
