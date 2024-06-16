package helpers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	StartupTimeoutSeconds = 5
	LogOccurrence         = 2
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	scripts, err := getSQLScripts(filepath.Join("../../", "migrations"))
	if err != nil {
		log.Fatalf("error getting sql scripts: %v", err)
	}
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(scripts...),
		postgres.WithDatabase("postges-test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(LogOccurrence).WithStartupTimeout(StartupTimeoutSeconds*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

func getSQLScripts(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var scripts []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			scripts = append(scripts, filepath.Join(dir, file.Name()))
		}
	}

	return scripts, nil
}
