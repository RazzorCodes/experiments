package main

import (
	"context"
	mcpdefs "razzor/golang-mcp/internal/app/mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	logger "razzor/golang-mcp/internal/utils"
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

	tm, err := mcpdefs.NewTrilliumMcp()
	if err != nil {
		logger.Get().Fatal("Failed to initialize: " + err.Error())
	}

	ctx := context.Background()
	transport := &mcp.StdioTransport{}
	err = tm.Server.Run(ctx, transport)
	if err != nil {
		logger.Get().Fatal("Server error: " + err.Error())
	}
}
