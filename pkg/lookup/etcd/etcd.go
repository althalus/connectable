package etcd

import (
	"strings"
	"time"

	env "github.com/MattAitchison/envconfig"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/gliderlabs/connectable/pkg/lookup"
)

var (
	// envconfig seems to not handle stringlist right
	//etcdPeers  = env.StringList("etcd_peers", []string{"http://127.0.0.1:4001"}, "Etcd Peers")
	etcdPeers  = env.String("etcd_peers", "http://127.0.0.1:4001", "Comman separated list of etcd peers")
	etcdPrefix = env.String("etcd_prefix", "", "Prefix to look for services in.")
	peersList []string
)

func init() {
	lookup.Register("etcd", new(etcdResolver))
	peersList = strings.Split(etcdPeers, ",")
}

type etcdResolver struct{}

func (r *etcdResolver) Lookup(addr string) ([]string, error) {
	cfg := client.Config{
		Endpoints: peersList,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	path := etcdPrefix + "/" + addr
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}
	var addrs []string
	for _, node := range resp.Node.Nodes {
		addrs = append(addrs, node.Value)
	}
	return addrs, nil
}
