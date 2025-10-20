// Incompatible change of how environment variables are handled in ~v1.17.0
// Filed issue https://github.com/spf13/viper/issues/1717.
package config_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ViperConfig struct {
	Test    ViperTest
	TestPtr *ViperTest
}

type ViperTest struct {
	String string
}

func TestViper(t *testing.T) {
	// Given
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("test.string", "test-original")
	v.SetDefault("testptr.string", "test-original")

	// When
	t.Setenv("TEST.STRING", "test-changed")
	t.Setenv("TESTPTR.STRING", "test-changed")
	config := &ViperConfig{}
	err := v.Unmarshal(config)
	require.NoError(t, err)

	// Then
	assert.Equal(t, "test-changed", config.Test.String)
	assert.Equal(t, "test-changed", v.Get("test.string"))
	assert.Equal(t, "test-changed", config.TestPtr.String)
	assert.Equal(t, "test-changed", v.Get("testptr.string"))
}
