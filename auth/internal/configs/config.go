package configs

import (
	"strings"

	"github.com/spf13/viper"
)

// APP_NESTED_HOST="testhost" APP_NESTED_PORT="8080" APP_SAMPLE_FIELD="foo" make dev
type Config struct {
	Nested struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"nested"`
	SampleField string `mapstructure:"sample_field"`
}

func Load(appName string) (*Config, error) {
	defaultViper := viper.New()

	// Set default config here

	// Set config file here

	// Merge configs

	// Set environment variables here
	defaultViper.SetEnvPrefix("APP")
	defaultViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	defaultViper.AutomaticEnv()

	// Ensure Viper knows about the keys so env-only values are picked up.
	// Without these, AutomaticEnv will not create new keys for Unmarshal.
	_ = defaultViper.BindEnv("nested.host")
	_ = defaultViper.BindEnv("nested.port")
	_ = defaultViper.BindEnv("sample_field")

	var cfg Config
	if err := defaultViper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
