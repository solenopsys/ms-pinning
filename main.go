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

	api := Api{
		addr: ":" + os.Getenv("api.Port"),
	}

	api.Start()

}
