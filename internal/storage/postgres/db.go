package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"time"
)

type DB struct {
	DB       *sqlx.DB
	Conn     *pgx.Conn
	pool     *dockertest.Pool
	resource *dockertest.Resource
	dbURL    string
}

func NewDB() *DB {
	return &DB{}
}

func (d *DB) SetupDB(migrationPath string) error {
	var err error
	d.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not connect to docker: %w", err)
	}

	d.resource, err = d.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		return fmt.Errorf("could not start postgres container: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = d.pool.Retry(func() error {
		d.dbURL = fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", d.resource.GetPort("5432/tcp"))
		d.DB, err = sqlx.ConnectContext(ctx, "postgres", d.dbURL)
		if err != nil {
			return fmt.Errorf("could not connect to postgres: %w", err)
		}

		d.Conn, err = pgx.Connect(ctx, d.dbURL)
		if err != nil {
			return fmt.Errorf("could not connect to postgres: %w", err)
		}

		err = Migration(ctx, migrationPath, d.DB)
		if err != nil {
			return fmt.Errorf("could not migrate: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not connect to postgres container: %w", err)
	}
	return nil
}

func (d *DB) TeardownDB() {
	if err := d.pool.Purge(d.resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %v", err)
	}
	err := d.DB.Close()
	if err != nil {
		log.Fatalf("Could not close Docker resource: %v", err)
	}
}
