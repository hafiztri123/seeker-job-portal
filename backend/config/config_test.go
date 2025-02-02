package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("LoadConfig", func(t *testing.T) {

		config, err := Load()
		require.NoError(t, err)

		assert.Equal(t, "hafizh", config.Database.User)
		assert.Equal(t, "Sudarmi12", config.Database.Password)
		assert.Equal(t, "172.31.236.79", config.Database.Host)
		assert.Equal(t, "5432", config.Database.Port)
		assert.Equal(t, "seeker", config.Database.Name)
	})

	t.Run("ValidateDatabaseConfig", func(t *testing.T) {
		os.Setenv("DB_PORT", "invalid")
		_, err := Load()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "database port is required")		
	})

	t.Run("GetDSN", func(t *testing.T) {
		os.Setenv("DB_USER", "user")
        os.Setenv("DB_PASSWORD", "pass")
        os.Setenv("DB_HOST", "localhost")
        os.Setenv("DB_PORT", "5432")
        os.Setenv("DB_NAME", "dbname")

		config, err := Load()
		require.NoError(t, err)

		expected := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
		assert.Equal(t, expected, config.Database.GetDSN())

	})

	





}