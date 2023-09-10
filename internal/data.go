package internal

import (
	"context"
	"errors"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"k8s.io/klog/v2"
	"ms-pinning/pkg"
	"strings"
	"time"
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

func changePattern(value string) string {
	newValue := strings.ReplaceAll(value, "*", "%")
	if newValue == "" {
		return "%"
	} else {
		return newValue
	}
}

func (d *Data) StatByTypes() (map[string]uint, error) {
	stat := make(map[string]uint)

	rows, err := d.Connection.Query(context.Background(), "select name, count(*) from pins join public.labels l on pins.id = l.pin_id group by name")

	for rows.Next() {
		var name string
		var count uint
		err := rows.Scan(&name, &count)
		if err != nil {
			return nil, err
		} else {
			stat[name] = count
		}
	}

	return stat, err
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
				klog.Info(i, " Pin saved in db: ", pinId)
				err := d.AddLabels(pinConf.Labels, pinId, userId, false)
				if err != nil {
					return err
				} else {
					klog.Info("Label save in db: ", pinId)
				}
			}
		}
	}
	return nil
}

func (d *Data) AddLabels(labels map[string]string, pinId string, userId uint64, checkUser bool) error {
	if checkUser {
		query := "select user_id from pins where id=$1 and  user_id=$2"
		var userIdFromDb uint64

		err := d.Connection.QueryRow(context.Background(), query, pinId, userId).Scan(&userIdFromDb)
		if err != nil {
			return err
		}
		if userIdFromDb != userId {
			return errors.New("Mistake pin owner ")
		}
	}

	for name, value := range labels {
		err := d.AddLabel(name, value, pinId)
		if err != nil {
			return errors.New("Label save in db error: " + err.Error())
		} else {
			klog.Info("Label save in db: ", pinId, name)
		}
	}

	return nil
}

func (d *Data) DeleteLabels(pinId string, userId uint64) error {
	query := "delete from labels where labels.pin_id=$1 and labels.pin_id in (select id from pins where user_id=$2)"
	_, err := d.Connection.Exec(context.Background(), query, pinId, userId)

	return err
}

func convertRecordRow(rows pgx.Rows) (map[string]map[string]string, error) {
	defer rows.Close()
	result := make(map[string]map[string]string)

	for rows.Next() {
		var ipnsId string
		var name string
		var value string
		err := rows.Scan(&ipnsId, &name, &value)
		if err != nil {
			return nil, err
		} else {
			if result[ipnsId] == nil {
				result[ipnsId] = make(map[string]string)
			}
			result[ipnsId][name] = value
		}
	}
	return result, nil
}

func (d *Data) SelectIpns(namePattern string, valuePattern string) (map[string]map[string]string, error) {
	nameFilter := changePattern(namePattern)
	valueFilter := changePattern(valuePattern)

	query := "select (select id from ipns where ipns.pin_id=labels.pin_id) ipnsId,name,value from labels where pin_id in  (    select labels.pin_id  from labels join ipns on labels.pin_id = ipns.pin_id where labels.name like $1 and value like $2)"
	rows, err := d.Connection.Query(context.Background(), query, nameFilter, valueFilter)

	if err != nil {
		return nil, err
	}

	return convertRecordRow(rows)
}

func (d *Data) SelectPins(namePattern string, valuePattern string) (map[string]map[string]string, error) {
	nameFilter := changePattern(namePattern)
	valueFilter := changePattern(valuePattern)

	query := "select pin_id,name,value from labels where pin_id in (select pin_id from labels where name like $1 and value like $2)"
	rows, err := d.Connection.Query(context.Background(), query, nameFilter, valueFilter)

	if err != nil {
		return nil, err
	}

	return convertRecordRow(rows)
}

type RecordInfo struct {
	Id       string            `json:"cid"`
	CreateAt time.Time         `json:"create"`
	Labels   map[string]string `json:"labels"`
}

func (d *Data) GetPin(cid string) (*RecordInfo, error) {

	var pinId string
	var createAt time.Time
	err := d.Connection.QueryRow(context.Background(), "select id,created_at from pins where id = $1", cid).Scan(&pinId, &createAt)
	if err != nil {
		return nil, err
	}

	query := "select pin_id,name,value from labels where pin_id = $1"
	rows, err := d.Connection.Query(context.Background(), query, pinId)

	if err != nil {
		return nil, err
	}

	result, err := convertRecordRow(rows)

	if err != nil {
		return nil, err
	}

	return &RecordInfo{
		Id:       pinId,
		CreateAt: createAt,
		Labels:   result[cid],
	}, nil
}

func (d *Data) GetIpnsByName(name string, userId uint64) (*RecordInfo, error) {
	var ipnsId string
	var createAt time.Time
	err := d.Connection.QueryRow(context.Background(),
		"select id,created_at  from ipns where name = $1 and user_id=$2", name, userId).Scan(&ipnsId, &createAt)
	if err != nil {
		return nil, err
	}

	labels, err := d.GetIpnsLabels(ipnsId)
	if err != nil {
		return nil, err
	}

	return &RecordInfo{
		Id:     ipnsId,
		Labels: labels[ipnsId],
	}, nil
}

func (d *Data) GetIpnsById(cid string) (*RecordInfo, error) {
	var ipnsId string
	var createAt time.Time
	err := d.Connection.QueryRow(context.Background(),
		"select id,created_at  from ipns where id = $1", cid).Scan(&ipnsId, &createAt)
	if err != nil {
		return nil, err
	}

	labels, err := d.GetIpnsLabels(ipnsId)
	if err != nil {
		return nil, err
	}

	return &RecordInfo{
		Id:     ipnsId,
		Labels: labels[ipnsId],
	}, nil
}

func (d *Data) GetIpnsLabels(ipnsId string) (map[string]map[string]string, error) {

	query := " select ipns.id,labels.name,value  from labels join ipns on labels.pin_id = ipns.pin_id where ipns.id = $1 "
	rows, err := d.Connection.Query(context.Background(), query, ipnsId)

	if err != nil {
		return nil, err
	}

	return convertRecordRow(rows)
}

func (d *Data) PinExists(cid string) (bool, error) {
	var exists bool
	err := d.Connection.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pins WHERE id = $1)", cid).Scan(&exists)
	return exists, err
}

func (d *Data) NameExists(name string, id uint64) (bool, error) {
	var exists bool
	err := d.Connection.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM ipns WHERE name = $1 and user_id=$2)", name, id).Scan(&exists)
	return exists, err
}

func (d *Data) CheckNicNameExists(name string) (bool, error) {
	var exists bool
	err := d.Connection.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE nickname = $1)", name).Scan(&exists)
	return exists, err
}

func (d *Data) SetNicName(userId uint64, name string) (bool, error) {
	var exists bool
	err := d.Connection.QueryRow(context.Background(), "UPDATE users SET nickname=$1 WHERE id=$2", name, userId).Scan(&exists)
	return exists, err
}
