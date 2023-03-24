package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
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
