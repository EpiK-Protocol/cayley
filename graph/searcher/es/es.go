package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cayleygraph/quad"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/epik-protocol/epik-gateway-backend/clog"
	"github.com/epik-protocol/epik-gateway-backend/graph"
)

const (
	Name = "elasticsearch"

	IndexNode = "node"
	IndexMeta = "meta"
	DocIDMeta = "1"
)

func init() {
	graph.RegisterSearcher(Name, NewSearcher)
}

func NewSearcher(qs graph.QuadStore, opts graph.Options) (graph.Searcher, error) {
	addresses, err := opts.StringSliceKey("addresses", nil)
	if err != nil {
		return nil, err
	}
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	//
	if err = checkConnection(es); err != nil {
		return nil, err
	}
	//
	if err = CheckAndCreateIndex(es, IndexNode, index_node_mapping); err != nil {
		return nil, err
	}
	if err = CheckAndCreateIndex(es, IndexMeta, index_meta_mapping); err != nil {
		return nil, err
	}
	s := &elasticSeacher{
		es:   es,
		qs:   qs,
		quit: make(chan struct{}),
	}
	go s.syncStore()
	return s, nil
}

type elasticSeacher struct {
	quit chan struct{}
	es   *elasticsearch.Client
	qs   graph.QuadStore
}

func (s *elasticSeacher) syncStore() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-s.quit:
			clog.Errorf("[Elasticsearch] exit")
			return
		case <-ticker.C:
			clog.Errorf("[Elasticsearch] syncing run")
			s.qs.SyncToSearcher(ctx, s)
		}
	}
}

func (s *elasticSeacher) SearchNodes(keys map[string]struct{}) ([]quad.Value, error) {
	var outs []quad.Value
	for key := range keys {
		query := fmt.Sprintf(search_node, key)

		// Perform the search request.
		res, err := s.es.Search(
			s.es.Search.WithIndex(IndexNode),
			s.es.Search.WithBody(strings.NewReader(query)),
			s.es.Search.WithSize(1000), //TODO:
			// s.es.Search.WithDocvalueFields("content"),
			s.es.Search.WithPretty(),
		)
		if err != nil {
			clog.Errorf("Error getting response: %s, keywords are '%s'", err, key)
			return nil, err
		}
		defer res.Body.Close()

		srcs, err := parseRes(res)
		if err != nil {
			return nil, err
		}

		for _, src := range srcs {
			out, ok := quad.AsValue(src["content"])
			if !ok {
				clog.Errorf("Failed to convert %T to quad.Value", src["id"])
				return nil, fmt.Errorf("Failed to convert %T to quad.Value", src["id"])
			}
			outs = append(outs, out)
		}
	}
	return outs, nil
}

func (s *elasticSeacher) IndexNodes(nodes []graph.Node) error {
	for _, node := range nodes {
		// delete old document
		res, err := s.es.DeleteByQuery(
			[]string{IndexNode},
			strings.NewReader(fmt.Sprintf(`{
				"query":{
					"match":{
						"id": %d
					}
				}
			}`, node.ID)),
			s.es.DeleteByQuery.WithRefresh(true),
		)
		if err != nil {
			clog.Errorf("Error deleting by query: %s, node.id: %d", err, node.ID)
			return err
		}
		consumeRes(res, false, "delete old doc")
		res.Body.Close()

		// add new document
		data, err := json.Marshal(node)
		if err != nil {
			clog.Errorf("Error marshalling node: %s", err)
			return err
		}
		res, err = s.es.Index(
			IndexNode,
			bytes.NewReader(data),
		)
		if err != nil {
			clog.Errorf("Error indexing node: %s, node.id: %d", err, node.ID)
			return err
		}
		consumeRes(res, false, "index new doc")
		res.Body.Close()
	}
	return nil
}

func (s *elasticSeacher) SetMeta(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	res, err := s.es.Update(
		IndexMeta,
		DocIDMeta,
		strings.NewReader(fmt.Sprintf(`{
			"doc": %s
		}`, string(data))),
		s.es.Update.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return parseResError(res)
	}
	consumeRes(res, false, "")
	return nil
}

func (s *elasticSeacher) GetMeta(key string) (interface{}, error) {
	res, err := s.es.Get(
		IndexMeta,
		DocIDMeta,
		s.es.Get.WithPretty(),
	)
	if err != nil {
		clog.Errorf("Error getting meta: %s, key: %s", err, key)
		return "", err
	}
	defer res.Body.Close()

	src, err := parseGetRes(res)
	if err != nil {
		return "", err
	}
	return src[key], nil
}

func (s *elasticSeacher) Close() error {
	close(s.quit)
	return nil
}

func parseResError(res *esapi.Response) error {
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			clog.Errorf("Error parsing the response error body: %s", err)
			return err
		} else {
			status := e["status"].(float64)
			switch status {
			case 404:
				return ErrNotFound
			}
			mError := e["error"].(map[string]interface{})
			reason := mError["reason"].(string)

			return &EsError{
				Reason: reason,
				Status: status,
			}
		}
	}
	return nil
}

func parseRes(res *esapi.Response) ([]map[string]interface{}, error) {
	if res.IsError() {
		return nil, parseResError(res)
	}

	var (
		r       map[string]interface{}
		sources []map[string]interface{}
	)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		clog.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}

	clog.Infof(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		clog.Infof(" * _ID=%s, %s", hit.(map[string]interface{})["_id"], source)
		sources = append(sources, source.(map[string]interface{}))
	}
	return sources, nil
}

func parseGetRes(res *esapi.Response) (map[string]interface{}, error) {
	if res.IsError() {
		return nil, parseResError(res)
	}

	var (
		r map[string]interface{}
	)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		clog.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}
	return r["_source"].(map[string]interface{}), nil
}

func checkConnection(es *elasticsearch.Client) error {
	var r map[string]interface{}
	res, err := es.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		clog.Errorf("Error parsing the ES response body: %s", err)
		return err
	}
	clog.Infof("ES Client: %s", elasticsearch.Version)
	clog.Infof("ES Server: %s", r["version"].(map[string]interface{})["number"])

	io.Copy(ioutil.Discard, res.Body)
	return nil
}

func consumeRes(res *esapi.Response, log bool, logPrefix string) {
	if res == nil {
		return
	}
	if log {
		var buf bytes.Buffer
		io.Copy(&buf, res.Body)
		clog.Infof("%s: %s", logPrefix, string(buf.Bytes()))
		return
	}
	io.Copy(ioutil.Discard, res.Body)
}

func CheckAndCreateIndex(client *elasticsearch.Client, index, tpl string) error {
	res, err := client.Indices.Get([]string{index})
	if err != nil {
		clog.Errorf("Error getting index %s: %s", index, err)
		return err
	}
	defer res.Body.Close()

	if !res.IsError() {
		consumeRes(res, true, fmt.Sprintf("Index [%s] exists", index))
		return nil
	}

	err = parseResError(res)
	if err != ErrNotFound {
		return err
	}
	res2, err := client.Indices.Create(
		index,
		client.Indices.Create.WithBody(strings.NewReader(tpl)),
	)
	if err != nil {
		clog.Errorf("Error creating index %s: %s", index, err)
		return err
	}
	defer res2.Body.Close()

	if !res2.IsError() {
		consumeRes(res2, true, fmt.Sprintf("Index [%s] created", index))
		if index == IndexMeta {
			data, _ := json.Marshal(map[string]int{"maxid": 0})
			res, err = client.Index(
				IndexMeta,
				bytes.NewReader(data),
				client.Index.WithDocumentID(DocIDMeta),
			)
			if err != nil {
				return err
			}
			if res.IsError() {
				return parseResError(res)
			}
		}
		return nil
	}
	return parseResError(res2)
}
