package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// default configuraitons parameters
const (
	defaultServiceName = "groupie_tracker"
	// HTTP Server
	defaultHTTPServerHost         = "localhost"
	defaultHTTPServerPort         = "8000"
	defaultHTTPServerIdleTimeout  = time.Second * 30
	defaultHTTPServerWriteTimeout = time.Second * 10
	defaultHTTPServerReadTimeout  = time.Second * 10
	defaultMaxHeaderMb            = 3 << 20 // 3 mb
	// HTTP Client
	defaultHTTPClinetTimeout = time.Second * 15

	// Logger
	defaultLoggerLevel     = -4
	defaultLoggerSourceKey = true
	defaultLoggerOutput    = "stdout"
	defaultLoggerHandler   = "json"
)

// Config structures
type (
	Config struct {
		ServiceName string `json:"service_name"` // <- []byte(default.json) service_name
		HTTPServer  `json:"http_server"`
		HTTPClient  `json:"http_client"`
		// DBConfig
		Logger `json:"logger"` // encoding/json -> "logger" -> Logger "snake_case"
	}

	HTTPServer struct {
		Host         string        `json:"host"`
		Port         string        `json:"port"`
		IdleTimeout  time.Duration `json:"idle_timeout"`
		WriteTimeout time.Duration `json:"write_timeout"`
		ReadTimeout  time.Duration `json:"read_timeout"`
		MaxHeaderMb  int           `json:"max_header_mb"`
	}

	HTTPClient struct {
		Timeout time.Duration `json:"timeout"`
	}

	Logger struct {
		Level     int    `json:"level"`
		SourceKey bool   `json:"source_key"`
		Output    string `json:"stdout"`
		Handler   string `json:"handler"`
	}
)

const (
	configDir  = "configs"
	configFile = "default.json"
)

// InitConfig ...
func InitConfig() (*Config, error) {
	cfg := &Config{}

	// 1. setup config parameters from defautl configuraiton constants

	// 2. parse config file, read and validate configurations parameters
	//  and set them on config structure

	// 3. from environment variables

	// 4. validate configuraiotn parameters

	// 5. set this parameters on the Config structure

	return setupConfig(cfg)
}

// setupConfig ...
func setupConfig(cfg *Config) (config *Config, err error) {
	// parsing configuration file
	config, err = parseFileAndSetConfig(cfg)
	if err != nil {
		// populateDefaults
		return populateDefaults(cfg), err
	}

	return config, err

	// something
}

// parseConfigFile ...
func parseFileAndSetConfig(cfg *Config) (config *Config, err error) {
	filePath := strings.Join([]string{configDir, configFile}, "/")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// validate configurations parameters

	return cfg, err
}

// populateDefaults
func populateDefaults(cfg *Config) *Config {
	// Service
	cfg.ServiceName = defaultServiceName

	// HTTPServer
	cfg.HTTPServer.Host = defaultHTTPServerHost
	cfg.HTTPServer.Port = defaultHTTPServerPort
	cfg.HTTPServer.IdleTimeout = defaultHTTPServerIdleTimeout
	cfg.HTTPServer.WriteTimeout = defaultHTTPServerWriteTimeout
	cfg.HTTPServer.ReadTimeout = defaultHTTPServerReadTimeout
	cfg.HTTPServer.MaxHeaderMb = defaultMaxHeaderMb

	// HTTPClinet
	cfg.HTTPClient.Timeout = defaultHTTPClinetTimeout

	// Logger
	cfg.Logger.Level = defaultLoggerLevel
	cfg.Logger.SourceKey = defaultLoggerSourceKey
	cfg.Logger.Output = defaultLoggerOutput
	cfg.Logger.Handler = defaultLoggerHandler

	return cfg
}

//hello
