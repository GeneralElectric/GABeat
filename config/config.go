// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"
	"time"
)

type Config struct {
	Period          time.Duration `config:"period"`
	Googleanalytics GoogleAnalyticsConfig
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}

type GoogleAnalyticsConfig struct {
	GoogleCredentialsFilePath string `config:"google_credentials_file"`
	GoogleAnalyticsIDs        string `config:"ga_ids"`
	GoogleAnalyticsMetrics    string `config:"ga_metrics"`
	GoogleAnalyticsDimensions string `config:"ga_dimensions"`
}

func (gaConfig GoogleAnalyticsConfig) String() string {
	return fmt.Sprintf("GoogleCredentialsFilePath: %s, GoogleAnalyticsIDs: %s, GoogleAnalyticsMetrics: %s, GoogleAnalyticsDimensions %s",
		gaConfig.GoogleCredentialsFilePath,
		gaConfig.GoogleAnalyticsIDs,
		gaConfig.GoogleAnalyticsMetrics,
		gaConfig.GoogleAnalyticsDimensions)
}

func (gbConfig Config) String() string {
	return fmt.Sprintf("Period: %s, GoogleAnalyticsConfig: %s, ",
		gbConfig.Period,
		gbConfig.Googleanalytics)
}
