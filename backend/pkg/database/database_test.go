package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	t.Run("SuccessfulConnection", func(t *testing.T) {
		db, err := connect()
		assert.NoError(t, err)
		assert.NoError(t, db.Ping())
		defer db.Close()
	})

	t.Run("FailedConnection", func(t *testing.T) {
		db, err := connect()
		assert.Error(t, err)
		assert.Error(t, db.Ping())
	})

	t.Run("ConnectionPool" , func(t *testing.T) {
		db, err := connect()
		assert.NoError(t, err)
		assert.NoError(t, db.Ping())
		defer db.Close()

		stats := db.Stats()
		assert.Equal(t, 0, stats.InUse)
		assert.NotZero(t, stats.MaxOpenConnections)
	})


}