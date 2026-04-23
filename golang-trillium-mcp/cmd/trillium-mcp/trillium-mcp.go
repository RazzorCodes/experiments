package main

import (
	governor "razzor/trillium-mcp/internal/app"
	logger "razzor/trillium-mcp/internal/utils"
)

const (
	UserAgentName    = "trillium-mcp"
	UserAgentVersion = "0.0.1"
	UserAgent        = UserAgentName + "/" + UserAgentVersion
	// ---
	TrilliumAddress  = "http://notes.lan/etapi"
	TrilliumApiToken = "CBkpReV6TEM6_FumFTgWmj1N3BWVkRnYd9RKtPJVdhkGeRXpInKnXfns="
)

func main() {
	logger.Get().Info("=== mcp === trillium-mcp === golang ===")
	logger.Get().Info("mcp: " + UserAgent)
	logger.Get().Info("Initializing...")

	if governor.Setup() != nil {
		logger.Get().Error("Initialization failed!")
	} else {
		logger.Get().Info("Initialization completed")
	}
}
