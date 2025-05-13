package network

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSTAConfig(t *testing.T) {
	uuid := uuid.New()
	cfg := DefaultSTAConfig(uuid, "test", "test", true)
	assert.Equal(t, "802-11-wireless", cfg.Type)
	assert.Equal(t, "test", cfg.ID)
	assert.Equal(t, "test", cfg.SSID)
	assert.Equal(t, true, cfg.AutoConnect)
}

func TestDefaultAPConfig(t *testing.T) {
	uuid := uuid.New()
	cfg := DefaultAPConfig(uuid, "test", "test", true)
	assert.Equal(t, "802-11-wireless", cfg.Type)
	assert.Equal(t, "test", cfg.ID)
	assert.Equal(t, "test", cfg.SSID)
	assert.Equal(t, true, cfg.AutoConnect)
}
