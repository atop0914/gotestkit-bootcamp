package testcontainer

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// TestContainerOption tests the container helper functions
func TestContainerOption(t *testing.T) {
	t.Run("WithExposedPorts", func(t *testing.T) {
		opt := WithExposedPorts("5432/tcp", "8080/tcp")
		req := &testcontainers.ContainerRequest{}
		opt(req)
		if len(req.ExposedPorts) != 2 {
			t.Errorf("expected 2 exposed ports, got %d", len(req.ExposedPorts))
		}
	})

	t.Run("WithEnv", func(t *testing.T) {
		opt := WithEnv(map[string]string{"KEY": "value"})
		req := &testcontainers.ContainerRequest{}
		opt(req)
		if req.Env["KEY"] != "value" {
			t.Errorf("expected env KEY=value, got %s", req.Env["KEY"])
		}
	})
}

func TestConfigStruct(t *testing.T) {
	cfg := Config{
		Image:        "nginx:latest",
		ExposedPorts: []string{"80/tcp"},
		Env:          map[string]string{"NGINX_HOST": "localhost"},
	}

	if cfg.Image != "nginx:latest" {
		t.Errorf("expected image nginx:latest, got %s", cfg.Image)
	}
	if len(cfg.ExposedPorts) != 1 {
		t.Errorf("expected 1 exposed port, got %d", len(cfg.ExposedPorts))
	}
	if cfg.Env["NGINX_HOST"] != "localhost" {
		t.Errorf("expected NGINX_HOST=localhost, got %s", cfg.Env["NGINX_HOST"])
	}
}

func TestPostgresOptions(t *testing.T) {
	// Test that postgres options work correctly
	opts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	}
	
	_ = opts // just verify the options exist
}

func TestMySQLOptions(t *testing.T) {
	// Test that mysql options work correctly
	opts := []testcontainers.ContainerCustomizer{
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
	}
	
	_ = opts // just verify the options exist
}

func TestMongoOptions(t *testing.T) {
	// Test that mongo options work correctly  
	opts := []testcontainers.ContainerCustomizer{
		mongodb.WithUsername("root"),
		mongodb.WithPassword("root"),
	}
	
	_ = opts // just verify the options exist
}

func TestReadyTimeout(t *testing.T) {
	ctx := context.Background()
	err := Ready(ctx, 100*time.Millisecond)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestPostgresContainerConnectionString(t *testing.T) {
	// Test ConnectionString method exists and is callable
	var c *postgresContainer
	_ = c.ConnectionString
}

func TestMySQLContainerConnectionString(t *testing.T) {
	// Test ConnectionString method exists and is callable
	var c *mysqlContainer
	_ = c.ConnectionString
}

func TestRedisConnectionString(t *testing.T) {
	// Test that Redis container methods exist
	var c *redisContainer
	_ = c.ConnectionString
}

func TestMongoConnectionString(t *testing.T) {
	// Test that MongoDB container methods exist
	var c *mongoContainer
	_ = c.ConnectionString
}

func TestGenericContainerEndpoint(t *testing.T) {
	// Test that generic container methods exist
	var c *genericContainer
	_ = c.Endpoint
	_ = c.Host
}

// TestWithExposedPortsMultiple tests multiple ports
func TestWithExposedPortsMultiple(t *testing.T) {
	opt := WithExposedPorts("5432/tcp", "8080/tcp", "9090/udp")
	req := &testcontainers.ContainerRequest{}
	opt(req)
	if len(req.ExposedPorts) != 3 {
		t.Errorf("expected 3 ports, got %d", len(req.ExposedPorts))
	}
}

// TestWithEnvMultiple tests multiple env vars
func TestWithEnvMultiple(t *testing.T) {
	opt := WithEnv(map[string]string{"KEY1": "value1", "KEY2": "value2"})
	req := &testcontainers.ContainerRequest{}
	opt(req)
	if req.Env["KEY1"] != "value1" || req.Env["KEY2"] != "value2" {
		t.Error("expected both env vars to be set")
	}
}

// TestReadyContextCanceled tests context cancellation
func TestReadyContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := Ready(ctx, time.Second)
	if err == nil {
		t.Error("expected context error, got nil")
	}
}
