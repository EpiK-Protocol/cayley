package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/jsonld"
	"github.com/epik-protocol/gateway/graph"
	"github.com/epik-protocol/gateway/graph/memstore"
	chttp "github.com/epik-protocol/gateway/internal/http"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

var testData = []quad.Quad{
	{
		Subject:   quad.IRI("http://example.com/alice"),
		Predicate: quad.IRI("http://example.com/likes"),
		Object:    quad.IRI("http://example.com/bob"),
		Label:     nil,
	},
}

func serializeTestData() string {
	buffer := bytes.NewBuffer(nil)
	writer := jsonld.NewWriter(buffer)
	writer.WriteQuads(testData)
	writer.Close()
	return buffer.String()
}

func serve(addr string) {
	qs := memstore.New(testData...)
	qw, err := graph.NewQuadWriter("single", qs, graph.Options{})
	if err != nil {
		panic(err)
	}
	h := &graph.Handle{QuadStore: qs, QuadWriter: qw}
	chttp.SetupRoutes(h, &chttp.Config{})
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func TestGatewayExport(t *testing.T) {
	port, err := freeport.GetFreePort()
	require.NoError(t, err)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	uri := fmt.Sprintf("http://%s", addr)
	go serve(addr)
	time.Sleep(3)
	cmd := NewCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{
		"--uri",
		uri,
	})
	err = cmd.Execute()
	require.NoError(t, err)
	data := serializeTestData()
	require.NotEmpty(t, data)
	require.Equal(t, data, b.String())
}
