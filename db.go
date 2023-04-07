package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
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
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (db *Db) Disconnect() error {
	return db.connection.Close()
}

func (db *Db) Migrate() {
	driver, err := postgres.WithInstance(db.connection, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mydatabase",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
}
