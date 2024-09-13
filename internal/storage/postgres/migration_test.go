package postgres_test

import (
	"context"
	"fmt"
	"github.com/andreevym/gophkeeper/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"os"
	"path/filepath"
)

func migrate(db *sqlx.DB) error {
	return filepath.Walk("../../../migrations", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}

			err = postgres.ApplyMigration(context.TODO(), db, string(bytes))
			if err != nil {
				return fmt.Errorf("apply migration: %w", err)
			}
		}

		return nil
	})
}
