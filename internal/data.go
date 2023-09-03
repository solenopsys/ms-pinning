package internal

import (
	"context"
	"errors"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/jackc/pgx/v4/pgxpool"
	"k8s.io/klog/v2"
	"ms-pinning/pkg"
)

type Data struct {
	Connection *pgxpool.Pool
}

type User struct {
	id        uint64
	publicKey string
	createdAt string
}

type Pin struct {
	id        int
	userId    uint64
	createdAt string
	size      int
	state     string
}

func (d *Data) GetPinsCountByUser(userId int64) (int, error) {
	var count int
	err := d.Connection.QueryRow(context.Background(), "SELECT count(*) FROM  pins  WHERE pins.user_id= $1", userId).Scan(&count)

	return count, err
}

func (d *Data) GetPinsCount() (int, error) {
	var count int
	err := d.Connection.QueryRow(context.Background(), "SELECT count(*) FROM  pins ").Scan(&count)

	return count, err
}

func (d *Data) GetIpnsCount() (int, error) {
	var count int
	err := d.Connection.QueryRow(context.Background(), "SELECT count(*) FROM  ipns ").Scan(&count)

	return count, err
}

func (d *Data) GetUsersCount() (int, error) {
	var count int
	err := d.Connection.QueryRow(context.Background(), "SELECT count(*) FROM  users ").Scan(&count)
	return count, err
}

func (d *Data) GetUserById(publicKey string) (*User, error) {
	res, err := d.Connection.Query(context.Background(), "SELECT id,public_key,created_at FROM  users  WHERE public_key = $1", publicKey)
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

func (d *Data) AddUser(publicKey string) (uint64, error) {
	lastInsertId := uint64(0)
	query := "INSERT INTO users (public_key) VALUES ($1)"

	err := d.Connection.QueryRow(context.Background(), query, publicKey).Scan(&lastInsertId)

	return lastInsertId, err
}

func (d *Data) AddPin(id string, userId uint64, size uint64) error {

	query := "INSERT INTO pins (id,user_id, size) VALUES ($1,$2, $3) ON CONFLICT (id) DO NOTHING"
	_, err := d.Connection.Exec(context.Background(), query, id, userId, size)
	return err
}

func (d *Data) CreateIpnsRecord(id string, userId uint64, pinId string, name string) error {
	query := "INSERT INTO ipns (id,user_id, pin_id,name) VALUES ($1,$2,$3,$4) "
	_, err := d.Connection.Exec(context.Background(), query, id, userId, pinId, name)
	return err
}

func (d *Data) ChangeIpnsRecord(name string, pinId string, userId uint64) (string, error) {
	query := "UPDATE ipns set pin_id=$1::text WHERE name=$2 and user_id=$3 RETURNING id;"
	var id string
	err := d.Connection.QueryRow(context.Background(), query, pinId, name, userId).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (d *Data) AddLabel(name string, value string, pinId string) error {

	query := "INSERT INTO labels (name,value, pin_id) VALUES ($1,$2, $3) ON CONFLICT (name,value, pin_id) DO NOTHING"
	_, err := d.Connection.Exec(context.Background(), query, name, value, pinId)
	return err

}

func (d *Data) SavePins(allPins []pkg.PinConf, group map[string]*api.Pin, userId uint64) error {
	for i, pinConf := range allPins {
		pin := group[pinConf.Cid]
		if pin == nil {
			return errors.New("Pin not success: " + pinConf.Cid)
		} else {
			pinId := pin.Cid.String()
			err := d.AddPin(pinId, userId, 0)

			if err != nil {
				return errors.New("Pin save in db error: " + err.Error())
			} else {
				klog.Info(i, " Pin save in db: ", pinId)
				for name, value := range pinConf.Labels {
					err := d.AddLabel(name, value, pinId)
					if err != nil {
						return errors.New("Label save in db error: " + err.Error())
					} else {
						klog.Info(i, "Label save in db: ", pinId, name)
					}
				}
			}
		}
	}
	return nil
}
