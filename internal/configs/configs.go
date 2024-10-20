package configs

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	// server config
	ListenAddress  string `mapstructure:"listen_address"`
	MetricsAddress string `mapstructure:"metrics_address"`
	QuotesFile     string `mapstructure:"quotes_file"`
	MaxRequests    int64  `mapstructure:"max_requests"`
	PowDifficulty  int    `mapstructure:"pow_difficulty"`

	// client config
	ServerAddress string `mapstructure:"server_address"`
	Concurrency   int    `mapstructure:"concurrency"`
}

func newViper() *viper.Viper {
	vp := viper.New()

	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vp.SetDefault("listen_address", ":9101")
	vp.SetDefault("metrics_address", ":9102")
	vp.SetDefault("quotes_file", "quotes.csv")
	vp.SetDefault("max_requests", 1000)
	vp.SetDefault("pow_difficulty", 22)

	vp.SetDefault("server_address", "127.0.0.1:9101")
	vp.SetDefault("concurrency", 2)

	return vp
}

func load(cfgFile string) (*viper.Viper, error) {
	const op = "configs.load"

	vp := newViper()

	if cfgFile == "" {
		log.Info().Msg("config file not specified, using defaults")
		return vp, nil
	}

	vp.SetConfigFile(cfgFile)

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%s: failed to read config file: %w", op, err)
	}

	return vp, nil
}

func NewConfig(cfgFile string) (*Config, error) {
	const op = "configs.NewConfig"

	vp, err := load(cfgFile)
	if err != nil {
		return nil, err
	}

	var config Config

	err = vp.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to unmarshal config file: %w", op, err)
	}

	return &config, nil
}
