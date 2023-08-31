package pkg

import (
	"context"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs-cluster/ipfs-cluster/api/rest/client"
	"k8s.io/klog/v2"
)

type IpfsPinsGroup struct {
	Pins   []PinConf `json:"pins"`
	RepMin int       `json:"rep_min"`
	RepMax int       `json:"rep_max"`
}

type PinConf struct {
	Cid    string            `json:"cid"`
	Labels map[string]string `json:"labels"`
}

type IpfsCluster struct {
	Host   string
	Port   string
	client client.Client
}

func (cluster *IpfsCluster) Connect() {
	var err error

	cfg := client.Config{Host: cluster.Host, Port: cluster.Port}
	cluster.client, err = client.NewDefaultClient(
		&cfg)
	if err != nil {

		klog.Error("Failed to create client:", err)
		return
	}
}

func (cluster *IpfsCluster) PinGroup(group IpfsPinsGroup) (map[string]*api.Pin, error) {

	pins := map[string]*api.Pin{}

	for _, cid := range group.Pins {
		decodeCid, err := api.DecodeCid(cid.Cid)
		if err != nil {
			klog.Error("Failed decode cid:", err)
			return pins, err
		}

		pin, err := cluster.client.Pin(context.Background(), decodeCid, api.PinOptions{
			ReplicationFactorMin: group.RepMin,
			ReplicationFactorMax: group.RepMax,
		})

		if err != nil {
			klog.Error("Failed to pin CID:", err)
			return pins, err
		}
		pins[cid.Cid] = &pin
	}

	return pins, nil
}
