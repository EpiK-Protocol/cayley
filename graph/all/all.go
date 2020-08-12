package all

import (
	// supported backends
	_ "github.com/epik-protocol/epik-gateway-backend/graph/kv/all"
	_ "github.com/epik-protocol/epik-gateway-backend/graph/memstore"
	_ "github.com/epik-protocol/epik-gateway-backend/graph/nosql/all"
	_ "github.com/epik-protocol/epik-gateway-backend/graph/sql/cockroach"
	_ "github.com/epik-protocol/epik-gateway-backend/graph/sql/mysql"
	_ "github.com/epik-protocol/epik-gateway-backend/graph/sql/postgres"
)
