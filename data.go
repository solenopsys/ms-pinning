package main

import "database/sql"

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
	res, err := d.connection.Query("SELECT * FROM users WHERE public_key = '" + publicKey + "'")
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
	query := "INSERT INTO users (public_key) VALUES ('" + publicKey + "')"
	res, err := d.connection.Exec(query)
	if err != nil {
		panic(err)
	}

	return res.LastInsertId()
}

func (d *Data) AddPin(id string, userId int64, size int) (int64, error) {
	query := "INSERT INTO pins (id,user_id, size) VALUES ('" + id + "'," + string(userId) + ", " + string(size) + ")"
	res, err := d.connection.Exec(query)
	if err != nil {
		panic(err)
	}
	return res.LastInsertId()
}
