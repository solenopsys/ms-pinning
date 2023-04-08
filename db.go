package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"k8s.io/klog/v2"
)

type Db struct {
	name       string
	password   string
	username   string
	host       string
	port       string
	connection *sql.DB
}

func (db *Db) Connect() error {
	connectString := "postgres://" + db.username + ":" + db.password + "@" + db.host + ":" + db.port + "/" + db.name + "?sslmode=disable"
	var err error
	db.connection, err = sql.Open("postgres", connectString)
	return err
}

func (db *Db) Disconnect() error {
	return db.connection.Close()
}

func (db *Db) Migrate() {
	driver, err := postgres.WithInstance(db.connection, &postgres.Config{})
	if err != nil {
		klog.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		db.name,
		driver,
	)
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			klog.Info("No migrations to apply")
			klog.Error(err.Error())
		} else {
			klog.Fatal(err)
		}
	} else {
		klog.Info("Migrations applied")
	}
}
