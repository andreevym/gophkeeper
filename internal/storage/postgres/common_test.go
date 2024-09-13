package postgres_test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"

	_ "github.com/lib/pq"
)

var testDB *sqlx.DB
var conn *pgx.Conn

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource
	dbURL    string
)

func setupDB() error {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not connect to docker: %w", err)
	}

	resource, err = pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=testuser",
		"POSTGRES_PASSWORD=testpass",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		return fmt.Errorf("could not start postgres container: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pool.Retry(func() error {
		dbURL = fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", resource.GetPort("5432/tcp"))
		testDB, err = sqlx.ConnectContext(ctx, "postgres", dbURL)
		if err != nil {
			return fmt.Errorf("could not connect to postgres: %w", err)
		}

		conn, err = pgx.Connect(ctx, dbURL)
		if err != nil {
			return fmt.Errorf("could not connect to postgres: %w", err)
		}

		err = migrate(testDB)
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

func teardownDB() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %v", err)
	}
	err := testDB.Close()
	if err != nil {
		log.Fatalf("Could not close Docker resource: %v", err)
	}
}

func TestMain(m *testing.M) {
	err := setupDB()
	if err != nil {
		log.Fatalf("Could not setup postgres container: %v", err)
	}
	code := m.Run()
	teardownDB()
	os.Exit(code)
}
