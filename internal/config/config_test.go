package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigs(t *testing.T) {

	t.Run("Should load proper connection string When env variables are set", func(t *testing.T) {
		t.Setenv(DBPassword, "password")
		t.Setenv(DBUser, "user")
		t.Setenv(DBHost, "local")
		t.Setenv(DBName, "not-tasks")
		t.Setenv(Port, ":5000")

		cfg := LoadConfig()

		assert.Equal(t, "user:password@tcp(local)/not-tasks?charset=utf8", cfg.DbConnection)
		assert.Equal(t, ":5000", cfg.AppPort)
	})

	t.Run("Should set default values for db connection When there is no host or db name", func(t *testing.T) {
		t.Setenv(DBPassword, "password")
		t.Setenv(DBUser, "user")

		cfg := LoadConfig()

		assert.Equal(t, "user:password@tcp(localhost)/tasks?charset=utf8", cfg.DbConnection)
		assert.Equal(t, "", cfg.AppPort)
	})

	t.Run("Should set docker host When is local container is set", func(t *testing.T) {
		t.Setenv(DBPassword, "password")
		t.Setenv(DBUser, "user")
		t.Setenv(IsLocalContainer, "true")

		cfg := LoadConfig()

		assert.Equal(t, "user:password@tcp(host.docker.internal)/tasks?charset=utf8", cfg.DbConnection)
		assert.Equal(t, "", cfg.AppPort)
	})
}
