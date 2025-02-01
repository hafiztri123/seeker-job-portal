package config

import (
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


}