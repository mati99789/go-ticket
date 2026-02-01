package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupDb(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()

	container, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	conntStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, conntStr)
	if err != nil {
		t.Fatalf("failed to create connection pool: %v", err)
	}

	migrationPath := "../../migrations"
	if err := runMigrations(ctx, t, pool, migrationPath); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Cleanup the container after the test
	t.Cleanup(func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres container: %v", err)
		}
	})

	return pool

}

func runMigrations(ctx context.Context, t *testing.T, pool *pgxpool.Pool, migrationPath string) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migration path %s does not exist", migrationPath)
	}

	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	var migration []os.DirEntry
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migration = append(migration, file)
		}
	}

	sort.Slice(migration, func(i, j int) bool {
		return migration[i].Name() < migration[j].Name()
	})

	if len(migration) == 0 {
		t.Logf("no migrations found in %s", migrationPath)
		return nil
	}

	for _, file := range migration {
		sql, err := os.ReadFile(filepath.Join(migrationPath, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}
		_, err = pool.Exec(ctx, string(sql))
		if err != nil {
			return fmt.Errorf("failed to execute migration file %s: %v", file.Name(), err)
		}

		t.Logf("migrated %s", file.Name())
	}

	return nil
}
