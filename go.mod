module github.com/epik-protocol/epik-gateway-backend

go 1.12

require (
	github.com/badgerodon/peg v0.0.0-20130729175151-9e5f7f4d07ca
	github.com/cayleygraph/quad v1.2.4
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/containerd/continuity v0.0.0-20190426062206-aaeac12a7ffc // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/cznic/mathutil v0.0.0-20170313102836-1447ad269d64
	github.com/d4l3k/messagediff v1.2.1 // indirect
	github.com/dennwc/graphql v0.0.0-20180603144102-12cfed44bc5d
	github.com/dgraph-io/badger v1.5.5 // indirect
	github.com/dlclark/regexp2 v1.1.4 // indirect
	github.com/docker/docker v0.7.3-0.20180412203414-a422774e593b // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dop251/goja v0.0.0-20190105122144-6d5bf35058fa
	github.com/elastic/go-elasticsearch/v7 v7.9.0
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-jsonrpc v0.1.1
	github.com/filecoin-project/specs-actors v0.0.0-00010101000000-000000000000
	github.com/flimzy/diff v0.1.6 // indirect
	github.com/fsouza/go-dockerclient v1.2.2
	github.com/go-sourcemap/sourcemap v2.1.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/hidal-go/hidalgo v0.0.0-20190814174001-42e03f3b5eaa
	github.com/ipfs/go-cid v0.0.6
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.5-0.20200428170625-a0bd04d3cbdf // indirect
	github.com/ipfs/go-ipld-format v0.2.0 // indirect
	github.com/ipfs/go-log/v2 v2.1.2-0.20200626104915-0016c0b4b3e4 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.3.0+incompatible
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/lib/pq v1.7.0
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/multiformats/go-multiaddr-net v0.1.5
	github.com/multiformats/go-multihash v0.0.14
	github.com/onsi/ginkgo v1.12.1 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/opencontainers/selinux v1.0.0 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/peterh/liner v0.0.0-20170317030525-88609521dc4b
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/piprate/json-gold v0.3.0
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.1.0 // indirect
	github.com/rogpeppe/go-internal v1.6.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.0
	github.com/tylertreat/BoomFilters v0.0.0-20181028192813-611b3dbe80e8
	github.com/warpfork/go-wish v0.0.0-20200122115046-b9ea61034e4a // indirect
	go.etcd.io/bbolt v1.3.4 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	golang.org/x/tools v0.0.0-20200729181040-64cdafbe085c // indirect
	google.golang.org/appengine v1.6.5
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/olivere/elastic.v5 v5.0.81 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace github.com/Sirupsen/logrus => github.com/Sirupsen/logrus v1.0.1

replace github.com/filecoin-project/specs-actors => ../specs-actors
