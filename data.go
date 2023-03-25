package main

type Data struct {
	db *Db
}

/*
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       public_key VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),

);
*/

type User struct {
	id        int64
	publicKey string
	createdAt string
}

/*
CREATE TABLE pins (

	id  PRIMARY KEY,
	user_id foreign key REFERENCES users(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	size bigint,
	state VARCHAR(255) NOT NULL DEFAULT 'new'

);
*/
type Pin struct {
	id        int
	userId    int64
	createdAt string
	size      int
	state     string
}

// function get user by id

func (data *Data) GetUserById(publicKey string) *User {
	res := data.db.Query("SELECT * FROM users WHERE public_key = '" + publicKey + "'")
	row := res.Next()
	if row {
		user := &User{}
		res.Scan(&user.id, &user.publicKey, &user.createdAt)
		return user
	} else {
		return nil
	}
}

// function add user to db
func (data *Data) AddUser(publicKey string) (int64, error) {
	insert, err := data.db.Insert("INSERT INTO users (public_key) VALUES ('" + publicKey + "')")
	if err != nil {
		panic(err)
	}
	id, err := insert.LastInsertId()
	return id, err
}

// function add pin to db
func (data *Data) AddPin(id string, userId int64, size int) {
	data.db.Insert("INSERT INTO pins (id,user_id, size) VALUES ('" + id + "'," + string(userId) + ", " + string(size) + ")")
}
