package pkg

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"k8s.io/klog/v2"
)

type Db struct {
	Name       string
	Password   string
	Username   string
	Host       string
	Port       string
	Connection *sql.DB
}

func (db *Db) Connect() error {
	connectString := "postgres://" + db.Username + ":" + db.Password + "@" + db.Host + ":" + db.Port + "/" + db.Name + "?sslmode=disable"
	var err error
	db.Connection, err = sql.Open("postgres", connectString)
	return err
}

func (db *Db) Disconnect() error {
	return db.Connection.Close()
}

func (db *Db) Migrate() {
	driver, err := postgres.WithInstance(db.Connection, &postgres.Config{})
	if err != nil {
		klog.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		db.Name,
		driver,
	)
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			klog.Info("No migrations to apply")
			klog.Error(err.Error())
		} else {
			klog.Fatal(err)
		}
	} else {
		klog.Info("Migrations applied")
	}
}
