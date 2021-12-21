package configs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	assert.Equal(t, &Config{}, NewConfig())
}

func TestConfig_Read(t *testing.T) {
	port := 4222
	exp := &Config{}
	exp.Nats.Port = port
	os.Setenv("FACEIT_SUBSCRIBER_NATS_PORT", fmt.Sprintf("%d", port))
	c := NewConfig()
	c.Read()
	assert.Equal(t, exp.Nats.Port, c.Nats.Port)
}
