package main

import (
	"context"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs-cluster/ipfs-cluster/api/rest/client"
	"k8s.io/klog/v2"
)

type IpfsPingGroup struct {
	Cids   []string `json:"cids"`
	RepMin int      `json:"rep_min"`
	RepMax int      `json:"rep_max"`
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

func (cluster *IpfsCluster) PinGroup(group IpfsPingGroup) []api.Pin {

	pins := []api.Pin{}

	for _, cid := range group.Cids {
		decodeCid, err := api.DecodeCid(cid)
		if err != nil {
			klog.Error("Failed decode cid:", err)
			return nil
		}

		pin, err := cluster.client.Pin(context.Background(), decodeCid, api.PinOptions{
			ReplicationFactorMin: group.RepMin,
			ReplicationFactorMax: group.RepMax,
		})

		pins = append(pins, pin)
		if err != nil {
			klog.Error("Failed to pin CID:", err)
			return nil
		}
	}

	return pins
}
