// Package testcontainer provides testcontainer integration helpers.
// This package provides convenience wrappers around testcontainers-go.
// Install testcontainers-go separately: go get github.com/testcontainers/testcontainers-go
package testcontainer

import (
	"context"
	"fmt"
	"time"
)

// Container defines a test container interface
type Container interface {
	Terminate(ctx context.Context) error
	Endpoint(ctx context.Context) (string, error)
}

// Config holds container configuration
type Config struct {
	Image        string
	ExposedPorts []string
	Env         map[string]string
}

// StartPostgres starts a PostgreSQL container for testing
func StartPostgres(ctx context.Context) (Container, error) {
	return startGeneric(ctx, Config{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
	})
}

// StartMySQL starts a MySQL container for testing
func StartMySQL(ctx context.Context) (Container, error) {
	return startGeneric(ctx, Config{
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "test",
			"MYSQL_PASSWORD":      "test",
		},
	})
}

// StartRedis starts a Redis container for testing
func StartRedis(ctx context.Context) (Container, error) {
	return startGeneric(ctx, Config{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		Env:         map[string]string{},
	})
}

// StartMongo starts a MongoDB container for testing
func StartMongo(ctx context.Context) (Container, error) {
	return startGeneric(ctx, Config{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
		},
	})
}

// startGeneric is a generic container starter
func startGeneric(ctx context.Context, cfg Config) (Container, error) {
	// Simulated container for demonstration
	// In real usage, integrate with testcontainers-go:
	// req := testcontainers.ContainerRequest{
	//     Image:        cfg.Image,
	//     ExposedPorts: cfg.ExposedPorts,
	//     Env:          cfg.Env,
	// }
	// c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	//     ContainerRequest: req,
	//     Started:          true,
	// })
	time.Sleep(100 * time.Millisecond) // simulate startup
	return &genericContainer{config: cfg}, nil
}

type genericContainer struct {
	config Config
	host   string
	port   string
}

func (c *genericContainer) Terminate(ctx context.Context) error {
	return nil
}

func (c *genericContainer) Endpoint(ctx context.Context) (string, error) {
	// In real implementation, use c.Host(ctx) and c.MappedPort(ctx, ...)
	return fmt.Sprintf("localhost:5432"), nil
}
