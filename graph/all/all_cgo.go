//+build cgo

package all

import (
	// backends requiring cgo
	_ "github.com/epik-protocol/gateway/graph/sql/sqlite"
)
