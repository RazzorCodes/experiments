package config

import (
	"fmt"
	logger "razzor/trillium-mcp/internal/utils"

	"github.com/joho/godotenv"
)

const (
	app_version_env_key = "GOLANG_TRILLIUM_MCP_VERSION"
	default_app_version = "0.0.0-error"

	etapi_address_env_key = "TRILLIUM_ETAPI_ADDRESS"
	default_etapi_address = "localhost:8000"
	etapi_apikey_env_key  = "TRILLIUM_ETAPI_APIKEY"
	default_etapi_apikey  = ""

	// default_etapi_user = "user"
	// default_etapi_pass = "pass"
)

type Config struct {
	AppVersion string

	EtapiAddress string
	EtapiApikey  string
}

func getEnv() (map[string]string, error) {
	env, err := godotenv.Read()
	if err != nil {
		logger.Error(fmt.Sprintf("Could not retrieve .env file: %s", err.Error()))
		return nil, err
	}

	return env, nil
}

func newConfig(env map[string]string) Config {
	result := Config{
		AppVersion:   default_app_version,
		EtapiAddress: default_etapi_address,
		EtapiApikey:  default_etapi_apikey,
	}

	if val, ok := env[app_version_env_key]; ok {
		result.AppVersion = val
	}

	if val, ok := env[etapi_address_env_key]; ok {
		result.EtapiAddress = val
	}

	if val, ok := env[etapi_apikey_env_key]; ok {
		result.EtapiApikey = val
	}

	return result
}

func GetConfig() (Config, error) {
	env, err := getEnv()
	if err != nil {
		logger.Warning("Governor could not retreive env. Falling back to defaults")
	}

	return newConfig(env), err
}
