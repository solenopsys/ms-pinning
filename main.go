package main

import (
	"k8s.io/klog/v2"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	db := Db{
		name:     "ms_pinning",
		password: os.Getenv("postgres.Password"),
		username: os.Getenv("postgres.User"),
		host:     os.Getenv("postgres.Host"),
		port:     os.Getenv("postgres.Port"),
	}
	ipfsClusterPort := os.Getenv("ipfs-cluster.Port")
	ipfsClusterHost := os.Getenv("ipfs-cluster.Host")
	apiPort := os.Getenv("api.Port")

	err := db.Connect()
	if err != nil {
		klog.Fatal(err)
	}

	db.Migrate()

	defer func(db *Db) {
		err := db.Disconnect()
		if err != nil {
			klog.Fatal(err)
		}
	}(&db)

	ipfsCluster := &IpfsCluster{
		Host: ipfsClusterHost, Port: ipfsClusterPort,
	}
	ipfsCluster.Connect()
	dataService := &Data{connection: db.connection}

	api := Api{
		addr: ":" + apiPort,
		ipfs: ipfsCluster,
		data: dataService,
	}

	api.Start()

}
