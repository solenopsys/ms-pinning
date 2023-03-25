package main

import (
	"context"
	"fmt"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs-cluster/ipfs-cluster/api/rest/client"
)

type IpfsPingGroup struct {
	Cids   []string `json:"cids"`
	RepMin int      `json:"rep_min"`
	RepMax int      `json:"rep_max"`
}

type IpfsCluster struct {
	IpfsClusterUrl string
	conection      client.Client
}

func (cluster *IpfsCluster) Connect() {
	var err error

	cfg := client.Config{Host: "bla"}
	cluster.conection, err = client.NewDefaultClient(
		&cfg)
	if err != nil {
		fmt.Println("Failed to create client:", err)
		return
	}
}

func (cluster *IpfsCluster) PinGroup(group IpfsPingGroup) []api.Pin {

	pins := []api.Pin{}

	for _, cid := range group.Cids {
		decodeCid, err := api.DecodeCid(cid)
		if err != nil {
			fmt.Println("Failed decode cid:", err)
			return nil
		}

		pin, err := cluster.conection.Pin(context.Background(), decodeCid, api.PinOptions{
			ReplicationFactorMin: group.RepMin,
			ReplicationFactorMax: group.RepMax,
		})

		pins = append(pins, pin)
		if err != nil {
			fmt.Println("Failed to pin CID:", err)
			return nil
		}
	}

	return pins
}
