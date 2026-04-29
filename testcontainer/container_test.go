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
