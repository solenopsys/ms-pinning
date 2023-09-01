package pkg

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"k8s.io/klog/v2"
	"time"
)

type Db struct {
	Name     string
	Password string
	Username string
	Host     string
	Port     string
	Pool     *pgxpool.Pool
}

func (db *Db) Connect() error {

	config, err := pgxpool.ParseConfig("user=" + db.Username + " host=" + db.Host + " port=" + db.Port + " password=" + db.Password + " dbname=" + db.Name + " sslmode=disable")

	if err != nil {
		klog.Fatal(err)
	}

	// Set connection pool options (optional)
	config.MaxConns = 50                       // Maximum number of connections in the pool
	config.MaxConnLifetime = 5 * time.Minute   // Maximum lifetime of a connection
	config.HealthCheckPeriod = 5 * time.Minute // Periodic health checks

	db.Pool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		klog.Fatal(err)
	}
	return err
}

func (db *Db) Disconnect() {
	db.Pool.Close()
}

func (db *Db) Migrate() {
	connectString := "postgres://" + db.Username + ":" + db.Password + "@" + db.Host + ":" + db.Port + "/" + db.Name + "?sslmode=disable"

	conn, err := sql.Open("postgres", connectString)
	defer conn.Close()

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
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
