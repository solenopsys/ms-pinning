package main

import (
	"flag"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"k8s.io/klog/v2"
	"ms-pinning/internal"
	"ms-pinning/pkg"
	"os"
)

const DEV_MODE = "dev"

var Mode string

func init() {
	flag.StringVar(&Mode, "mode", "", "a string var")
}

func main() {
	flag.Parse()
	devMode := Mode == DEV_MODE

	if devMode {
		err := godotenv.Load("info/.env")
		if err != nil {
			klog.Fatal("Error loading .env file")
		}
	}

	db := pkg.Db{
		Name:     "ms_pinning",
		Password: os.Getenv("postgres.Password"),
		Username: os.Getenv("postgres.User"),
		Host:     os.Getenv("postgres.Host"),
		Port:     os.Getenv("postgres.Port"),
	}

	ipfsClusterPort := os.Getenv("ipfs-cluster.Port")
	ipfsClusterHost := os.Getenv("ipfs-cluster.Host")

	apiPort := os.Getenv("api.Port")

	ipfsHost := os.Getenv("ipfs.Host")
	ipfsPort := os.Getenv("ipfs.Port")

	err := db.Connect()
	if err != nil {
		klog.Fatal(err)
	}

	db.Migrate()

	defer db.Disconnect()

	ipfsCluster := &pkg.IpfsCluster{
		Host: ipfsClusterHost, Port: ipfsClusterPort,
	}
	ipfsCluster.Connect()
	dataService := &internal.Data{Connection: db.Pool}

	api := internal.Api{
		Addr:        ":" + apiPort,
		IpfsCluster: ipfsCluster,
		Data:        dataService,
		Ipfs:        &pkg.IpfsNode{Host: ipfsHost, Port: ipfsPort},
	}

	api.Start()

}
