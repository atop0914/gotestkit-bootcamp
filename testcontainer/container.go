// Package testcontainer provides testcontainer integration helpers.
// This package provides convenience wrappers around testcontainers-go.
package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// Container defines a test container interface
type Container interface {
	Terminate(ctx context.Context, opts ...testcontainers.TerminateOption) error
	Endpoint(ctx context.Context) (string, error)
	Host(ctx context.Context) (string, error)
}

// Config holds container configuration
type Config struct {
	Image        string
	ExposedPorts []string
	Env          map[string]string
}

// ContainerOption configures a generic container request
type ContainerOption func(*testcontainers.ContainerRequest)

// WithExposedPorts sets exposed ports
func WithExposedPorts(ports ...string) ContainerOption {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = ports
	}
}

// WithEnv sets environment variables
func WithEnv(env map[string]string) ContainerOption {
	return func(req *testcontainers.ContainerRequest) {
		if req.Env == nil {
			req.Env = make(map[string]string)
		}
		for k, v := range env {
			req.Env[k] = v
		}
	}
}

// PostgresOption configures PostgreSQL container using functional options
type PostgresOption func(*postgres.PostgresContainer)

// WithPostgresDatabase sets the database name
func WithPostgresDatabase(dbName string) testcontainers.ContainerCustomizer {
	return postgres.WithDatabase(dbName)
}

// WithPostgresUsername sets the username
func WithPostgresUsername(user string) testcontainers.ContainerCustomizer {
	return postgres.WithUsername(user)
}

// WithPostgresPassword sets the password
func WithPostgresPassword(password string) testcontainers.ContainerCustomizer {
	return postgres.WithPassword(password)
}

// StartPostgres starts a PostgreSQL container for testing
func StartPostgres(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (Container, error) {
	c, err := postgres.Run(ctx, "postgres:15-alpine", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}
	return &postgresContainer{c}, nil
}

type postgresContainer struct {
	*postgres.PostgresContainer
}

func (c *postgresContainer) Host(ctx context.Context) (string, error) {
	return c.PostgresContainer.Host(ctx)
}

func (c *postgresContainer) Endpoint(ctx context.Context) (string, error) {
	port, err := c.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

// ConnectionString returns the full PostgreSQL connection string
func (c *postgresContainer) ConnectionString(ctx context.Context) (string, error) {
	return c.PostgresContainer.ConnectionString(ctx)
}

// MySQLOption configures MySQL container using functional options
type MySQLOption func(*mysql.MySQLContainer)

// WithMySQLDatabase sets the database name
func WithMySQLDatabase(dbName string) testcontainers.ContainerCustomizer {
	return mysql.WithDatabase(dbName)
}

// WithMySQLUsername sets the username
func WithMySQLUsername(user string) testcontainers.ContainerCustomizer {
	return mysql.WithUsername(user)
}

// WithMySQLPassword sets the password
func WithMySQLPassword(password string) testcontainers.ContainerCustomizer {
	return mysql.WithPassword(password)
}

// StartMySQL starts a MySQL container for testing
func StartMySQL(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (Container, error) {
	c, err := mysql.Run(ctx, "mysql:8", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start mysql: %w", err)
	}
	return &mysqlContainer{c}, nil
}

type mysqlContainer struct {
	*mysql.MySQLContainer
}

func (c *mysqlContainer) Host(ctx context.Context) (string, error) {
	return c.MySQLContainer.Host(ctx)
}

func (c *mysqlContainer) Endpoint(ctx context.Context) (string, error) {
	port, err := c.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

// ConnectionString returns the full MySQL connection string
func (c *mysqlContainer) ConnectionString(ctx context.Context) (string, error) {
	return c.MySQLContainer.ConnectionString(ctx)
}

// RedisOption configures Redis container using functional options
type RedisOption func(*redis.RedisContainer)

// StartRedis starts a Redis container for testing
func StartRedis(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (Container, error) {
	c, err := redis.Run(ctx, "redis:7-alpine", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start redis: %w", err)
	}
	return &redisContainer{c}, nil
}

type redisContainer struct {
	*redis.RedisContainer
}

func (c *redisContainer) Host(ctx context.Context) (string, error) {
	return c.RedisContainer.Host(ctx)
}

func (c *redisContainer) Endpoint(ctx context.Context) (string, error) {
	port, err := c.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

// ConnectionString returns the Redis connection string
func (c *redisContainer) ConnectionString(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := c.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

// MongoOption configures MongoDB container using functional options
type MongoOption func(*mongodb.MongoDBContainer)

// WithMongoUsername sets the username
func WithMongoUsername(user string) testcontainers.ContainerCustomizer {
	return mongodb.WithUsername(user)
}

// WithMongoPassword sets the password
func WithMongoPassword(password string) testcontainers.ContainerCustomizer {
	return mongodb.WithPassword(password)
}

// StartMongo starts a MongoDB container for testing
func StartMongo(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (Container, error) {
	c, err := mongodb.Run(ctx, "mongo:6", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start mongodb: %w", err)
	}
	return &mongoContainer{c}, nil
}

type mongoContainer struct {
	*mongodb.MongoDBContainer
}

func (c *mongoContainer) Host(ctx context.Context) (string, error) {
	return c.MongoDBContainer.Host(ctx)
}

func (c *mongoContainer) Endpoint(ctx context.Context) (string, error) {
	port, err := c.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

// ConnectionString returns the MongoDB connection string
func (c *mongoContainer) ConnectionString(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := c.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("mongodb://root:root@%s:%s", host, port.Port()), nil
}

// StartGeneric starts a generic container for testing
func StartGeneric(ctx context.Context, cfg Config, opts ...ContainerOption) (Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        cfg.Image,
		ExposedPorts: cfg.ExposedPorts,
		Env:          cfg.Env,
	}

	for _, opt := range opts {
		opt(&req)
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start generic container: %w", err)
	}

	return &genericContainer{c}, nil
}

type genericContainer struct {
	testcontainers.Container
}

func (c *genericContainer) Host(ctx context.Context) (string, error) {
	return c.Container.Host(ctx)
}

func (c *genericContainer) Endpoint(ctx context.Context) (string, error) {
	ports, err := c.Ports(ctx)
	if err != nil {
		return "", err
	}
	if len(ports) == 0 {
		return "", fmt.Errorf("no ports exposed")
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	// Try to find first valid port binding
	for portProto := range ports {
		bindings := ports[portProto]
		if len(bindings) > 0 && bindings[0].HostPort != "" {
			return fmt.Sprintf("%s:%s", host, bindings[0].HostPort), nil
		}
	}
	return "", fmt.Errorf("could not determine endpoint")
}

// Ready waits for container to be ready with a timeout
func Ready(ctx context.Context, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for container to be ready")
	}
}
