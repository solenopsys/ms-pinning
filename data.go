package main

import (
	"database/sql"
)

type Data struct {
	connection *sql.DB
}

type User struct {
	id        int64
	publicKey string
	createdAt string
}

type Pin struct {
	id        int
	userId    int64
	createdAt string
	size      int
	state     string
}

func (d *Data) GetUserById(publicKey string) *User {
	res, err := d.connection.Query("SELECT * FROM users WHERE public_key = $1", publicKey)
	defer res.Close()

	if err != nil {
		panic(err)
	}

	row := res.Next()
	if row {
		user := &User{}
		res.Scan(&user.id, &user.publicKey, &user.createdAt)
		return user
	} else {
		return nil
	}
}

func (d *Data) AddUser(publicKey string) (int64, error) {
	lastInsertId := int64(0)
	query := "INSERT INTO users (public_key) VALUES ($1)"

	err := d.connection.QueryRow(query, publicKey).Scan(&lastInsertId)

	return lastInsertId, err
}

func (d *Data) AddPin(id string, userId int64, size uint64) error {

	query := "INSERT INTO pins (id,user_id, size) VALUES ($1,$2, $3)"
	err := d.connection.QueryRow(query, id, userId, size).Err()

	return err
}

func (d *Data) AddLabel(name string, value string, pinId string) error {

	query := "INSERT INTO labels (name,value, pin_id) VALUES ($1,$2, $3)"
	err := d.connection.QueryRow(query, name, value, pinId).Err()

	return err
}
