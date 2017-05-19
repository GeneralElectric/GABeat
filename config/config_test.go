// +build !integration

package config

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig2(t *testing.T) {
	// Tests with different params from config file
	absPath, err := filepath.Abs("../tests/files/")

	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	config := &Config{}

	// Reads second config file
	err = cfgfile.Read(config, absPath+"/config2.yml")
	assert.Nil(t, err)

	fmt.Println(fmt.Sprintf("test config: %s ", config))

	assert.Equal(t, 0*time.Minute, config.Period)
}

func TestToString(t *testing.T) {
	gaConfig := GoogleAnalyticsConfig{"/", "ids", "metric", "dimension"}
	beatConfig := Config{30 * time.Minute, gaConfig}

	actual := beatConfig.String()

	assert.Equal(t, "Period: 30m0s, GoogleAnalyticsConfig: GoogleCredentialsFilePath: /, GoogleAnalyticsIDs: ids, GoogleAnalyticsMetrics: metric, GoogleAnalyticsDimensions dimension, ",
		actual)
}
