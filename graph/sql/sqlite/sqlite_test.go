//+build cgo

package sqlite

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/epik-protocol/gateway/graph"
	"github.com/epik-protocol/gateway/graph/sql/sqltest"
)

func makeSqlite(t testing.TB) (string, graph.Options, func()) {
	tmpFile, err := ioutil.TempFile("", fmt.Sprintf("gateway_test_%s*", Type))
	if err != nil {
		t.Fatalf("Could not create working directory: %v", err)
	}
	return fmt.Sprintf("file:%s?_loc=UTC", tmpFile.Name()), nil, func() {
		os.RemoveAll(tmpFile.Name())
	}
}

var conf = &sqltest.Config{
	TimeRound: true,
	TimeInMcs: false,
}

func TestSqlite(t *testing.T) {
	sqltest.TestAll(t, Type, makeSqlite, conf)
}

func BenchmarkSqlite(t *testing.B) {
	sqltest.BenchmarkAll(t, Type, makeSqlite, conf)
}
