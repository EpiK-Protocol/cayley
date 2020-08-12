//+build cgo

package all

import (
	// backends requiring cgo
	_ "github.com/epik-protocol/epik-gateway-backend/graph/sql/sqlite"
)
