package core

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Setenv("PLUGIN_KEYCLOAK_URL", "http://test")
	libConfig = getConfig()
}

func TestLibConfig(t *testing.T) {
	assert.Equal(t, "http://test", libConfig.KeyCloak.Url)
	os.Unsetenv("PLUGIN_KEYCLOAK_URL")
}
