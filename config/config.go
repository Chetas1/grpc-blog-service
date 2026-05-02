// Package config loads runtime configuration from `config/config.yaml`
// and environment variables (via Viper's AutomaticEnv). Environment variables
// override file values; nested keys use `_` separators
// (e.g. `GRPCSERVER_PORT=9091`).
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is the top-level service configuration.
type Config struct {
	GrpcServer    Server
	GatewayServer Server
	GrpcClient    Client
}

// Server holds host/port/protocol for a network server.
type Server struct {
	Host     string
	Port     int
	Protocol string
}

// Client holds connection details for the gRPC client smoke-test binary.
type Client struct {
	ServerAddress string
}

// LoadConfig reads `config/config.yaml`, applies environment overrides,
// and returns the parsed Config. Returns a wrapped error if the file is
// missing or the YAML cannot be unmarshalled into Config.
func LoadConfig() (Config, error) {
	var cfg Config

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	return cfg, nil
}
