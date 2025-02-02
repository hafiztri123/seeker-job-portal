package database

import (
	"testing"
	"time"

	"github.com/hafiztri123/config"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	t.Run("SuccessfulConnection", func(t *testing.T) {

		cfg := &config.DatabaseConfig{
			Host: "172.31.236.79",
			Port: "5432",
			User: "hafizh",
			Password: "Sudarmi12",
			Name: "seeker",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnectionMaxLifetime: 5 * time.Minute,
		}

		db, err := Connect(cfg)
		assert.NoError(t, err)
		assert.NoError(t, db.Ping())
		defer db.Close()
	})

	t.Run("FailedConnection", func(t *testing.T) {

		cfg := &config.DatabaseConfig{
			Host: "invalid",
			Port: "invalid",
			User: "invalid",
			Password: "invalid",
			Name: "invalid",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnectionMaxLifetime: 5 * time.Minute,
		}
		db, err := Connect(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("ConnectionPool" , func(t *testing.T) {
		cfg := &config.DatabaseConfig{
			Host: "172.31.236.79",
			Port: "5432",
			User: "hafizh",
			Password: "Sudarmi12",
			Name: "seeker",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnectionMaxLifetime: 5 * time.Minute,
		}
	
		db, err := Connect(cfg)
		assert.NoError(t, err)
		assert.NoError(t, db.Ping())
		defer db.Close()

		stats := db.Stats()
		assert.Equal(t, 0, stats.InUse)
		assert.NotZero(t, stats.MaxOpenConnections)
	})


}