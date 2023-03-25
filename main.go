package main

import (
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	db := Db{
		name:     "mydatabase",
		password: "password",
		username: "username",
		host:     "localhost",
		port:     5432,
	}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *Db) {
		err := db.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}(&db)

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
