package mongo

import (
	"github.com/epik-protocol/gateway/graph"
	"github.com/hidal-go/hidalgo/legacy/nosql"
	"github.com/hidal-go/hidalgo/legacy/nosql/mongo"

	//import hidal-go first so the registration of the no sql stores occurs before quadstore iterates for registration
	gnosql "github.com/epik-protocol/gateway/graph/nosql"
)

const Type = mongo.Name

func Create(addr string, opt graph.Options) (nosql.Database, error) {
	return mongo.Dial(addr, gnosql.DefaultDBName, nosql.Options(opt))
}

func Open(addr string, opt graph.Options) (nosql.Database, error) {
	return mongo.Dial(addr, gnosql.DefaultDBName, nosql.Options(opt))
}
