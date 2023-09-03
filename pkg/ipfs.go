package pkg

import (
	"context"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"k8s.io/klog/v2"
	"time"
)

type IpfsNode struct {
	Host string
	Port string
}

func (node *IpfsNode) GetShell() *shell.Shell {
	shell := shell.NewShell(node.Host + ":" + node.Port)

	return shell
}

func (node *IpfsNode) CreateKey(keyName string) (string, error) {

	key, err := node.GetShell().KeyGen(context.Background(), keyName)
	if err != nil {
		fmt.Println("Error creating IPNS key:", err)
		return "", err
	}

	klog.Info("IPNS key '" + keyName + "' created with ID: " + key.Id)

	return key.Id, nil
}

func (node *IpfsNode) Publish(cid string, keyName string) error {
	var lifeTime time.Duration
	//24h
	lifeTime = 24 * time.Hour * 7
	var ttl time.Duration
	ttl = 5 * time.Minute
	details, err := node.GetShell().PublishWithDetails(cid, keyName, lifeTime, ttl, false)
	if err != nil {
		fmt.Println("Error publishing:", err)
		return err
	}

	klog.Info("Content published to IPNS with name: ", keyName)
	klog.Info("Publish response: ", details.Value)
	return nil
}
