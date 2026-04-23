package governor

import (
	"razzor/trillium-mcp/internal/config"
	"razzor/trillium-mcp/internal/connectors"
	logger "razzor/trillium-mcp/internal/utils"
)

func Setup() error {
	logger.Info("Setting up...")

	appconf, err := config.GetConfig()
	if err != nil {
		logger.Info("Config was not loaded sucessfully")
	}

	logger.Info("AppVersion: " + appconf.AppVersion)
	logger.Info("EtapiAddress: " + appconf.EtapiAddress)

	_, err = connectors.NewTrilliumConnector(appconf)

	return nil
}
