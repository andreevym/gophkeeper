package postgres

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
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

func (d *DB) SetupDB(ctx context.Context, migrationPath string) error {
	var err error
	d.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not connect to docker: %w", err)
	}

	d.resource, err = d.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return fmt.Errorf("could not start postgres container: %w", err)
	}

	time.Sleep(5 * time.Second)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = d.pool.Retry(func() error {
		d.dbURL = fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", getHostPort(d.resource, "5432/tcp"))
		d.DB, err = sqlx.ConnectContext(ctx, "postgres", d.dbURL)
		if err != nil {
			return fmt.Errorf("could not connect to postgres: %w", err)
		}

		return d.DB.Ping()
	})
	if err != nil {
		return fmt.Errorf("could not connect to postgres container: %w", err)
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
}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id)
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}
	return u.Hostname() + ":" + resource.GetPort(id)
}

func (d *DB) TeardownDB() {
	if err := d.pool.Purge(d.resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %v", err)
	}
	//err := d.DB.Close()
	//if err != nil {
	//	log.Fatalf("Could not close Docker resource: %v", err)
	//}
}
