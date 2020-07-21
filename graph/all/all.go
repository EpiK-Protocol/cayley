package all

import (
	// supported backends
	_ "github.com/epik-protocol/gateway/graph/kv/all"
	_ "github.com/epik-protocol/gateway/graph/memstore"
	_ "github.com/epik-protocol/gateway/graph/nosql/all"
	_ "github.com/epik-protocol/gateway/graph/sql/cockroach"
	_ "github.com/epik-protocol/gateway/graph/sql/mysql"
	_ "github.com/epik-protocol/gateway/graph/sql/postgres"
)
