package postgres

import (
	"context"
	"fmt"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func ApplyMigration(ctx context.Context, db *sqlx.DB, sql string) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := db.ExecContext(rCtx, sql)
	if err != nil {
		return fmt.Errorf("failed apply sql '%s': %w", sql, err)
	}

	return nil
}

func Migration(ctx context.Context, migrationPath string, db *sqlx.DB) error {
	err := filepath.Walk(migrationPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			logger.Logger().Info("apply migration", zap.String("path", path))
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			sql := string(bytes)
			err = ApplyMigration(ctx, db, sql)
			if err != nil {
				return fmt.Errorf("apply migration, path '%s', sql '%s': %w", path, sql, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("exec migration: %w", err)
	}
	return nil
}
