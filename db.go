package main

import (
	"database/sql"
	"log"
)

type Db struct {
	name       string
	password   string
	username   string
	host       string
	port       int
	connection *sql.DB
}

func (db *Db) Connect() error {
	connectString := "postgres://" + db.username + ":" + db.password + "@" + db.host + ":" + string(db.port) + "/" + db.name + "?sslmode=disable"
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

func (db *Db) Query(query string) *sql.Rows {
	return db.Query(query)
}

// insert function
func (db *Db) Insert(query string) (sql.Result, error) {
	return db.connection.Exec(query)
}
