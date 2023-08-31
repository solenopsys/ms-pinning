package internal

import (
	"database/sql"
	"errors"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"k8s.io/klog/v2"
	"ms-pinning/pkg"
)

type Data struct {
	Connection *sql.DB
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

func (d *Data) GetPinsCountByUser(userId int64) (int, error) {
	var count int
	err := d.Connection.QueryRow("SELECT count(*) FROM  pins  WHERE pins.user_id= $1", userId).Scan(&count)

	return count, err
}

func (d *Data) GetPinsCount() (int, error) {
	var count int
	err := d.Connection.QueryRow("SELECT count(*) FROM  pins ").Scan(&count)

	return count, err
}

func (d *Data) GetUsersCount() (int, error) {
	var count int
	err := d.Connection.QueryRow("SELECT count(*) FROM  users ").Scan(&count)
	return count, err
}

func (d *Data) GetUserById(publicKey string) (*User, error) {
	res, err := d.Connection.Query("SELECT id,public_key,created_at FROM  users  WHERE public_key = $1", publicKey)
	defer res.Close()

	if err != nil {
		return nil, err
	}

	row := res.Next()
	if row {
		user := &User{}
		res.Scan(&user.id, &user.publicKey, &user.createdAt)
		return user, nil
	} else {
		return nil, nil
	}
}

func (d *Data) AddUser(publicKey string) (int64, error) {
	lastInsertId := int64(0)
	query := "INSERT INTO users (public_key) VALUES ($1)"

	err := d.Connection.QueryRow(query, publicKey).Scan(&lastInsertId)

	return lastInsertId, err
}

func (d *Data) AddPin(id string, userId int64, size uint64) error {

	query := "INSERT INTO pins (id,user_id, size) VALUES ($1,$2, $3)"
	err := d.Connection.QueryRow(query, id, userId, size).Err()

	return err
}

func (d *Data) AddLabel(name string, value string, pinId string) error {

	query := "INSERT INTO labels (name,value, pin_id) VALUES ($1,$2, $3)"
	err := d.Connection.QueryRow(query, name, value, pinId).Err()

	return err
}

func (d *Data) SavePins(allPins []pkg.PinConf, group map[string]*api.Pin, userId int64) error {
	for _, pinConf := range allPins {
		pin := group[pinConf.Cid]
		if pin == nil {
			return errors.New("Pin not success: " + pinConf.Cid)
		} else {
			pinId := pin.Cid.String()
			err := d.AddPin(pinId, userId, 0)

			if err != nil {
				return errors.New("Pin save in db error: " + err.Error())
			} else {
				klog.Info("Pin save in db: ", pinId)
				for name, value := range pinConf.Labels {
					err := d.AddLabel(name, value, pinId)
					if err != nil {
						return errors.New("Label save in db error: " + err.Error())
					} else {
						klog.Info("Label save in db: ", pinId, name)
					}
				}
			}
		}
	}
	return nil
}
