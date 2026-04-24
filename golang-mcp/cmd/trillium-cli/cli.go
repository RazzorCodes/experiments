package main

import (
	"context"
	"log"
	"os"
	clicommands "razzor/golang-mcp/internal/app/cli"
	"razzor/golang-mcp/internal/config"
	connectors "razzor/golang-mcp/internal/connector"
	logger "razzor/golang-mcp/internal/utils"
	"strings"

	"github.com/urfave/cli/v3" // imports as package "cli"
)

type TrilliumCli struct {
	conn *connectors.TrilliumConnector
}

func NewTrilliumCli() (TrilliumCli, error) {
	logger.Info("Setting up...")

	appconf, err := config.GetConfig()
	if err != nil {
		logger.Info("Config was not loaded sucessfully")
	}

	logger.Info("AppVersion: " + appconf.AppVersion)
	logger.Info("EtapiAddress: " + appconf.EtapiAddress)

	var t TrilliumCli
	t.conn, err = connectors.NewTrilliumConnector(appconf)
	if err != nil {
		logger.Fatal("")
	}

	return t, nil
}

func main() {
	logger.Info("=======================================")
	logger.Info("=== cli === trillium-cli === golang ===")

	trilliumCli, _ := NewTrilliumCli()

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "search for a keyword",
				Action:  clicommands.GetSearchAction(trilliumCli.conn),
			},
			{
				Name:    "content",
				Aliases: []string{"c"},
				Usage:   "get content of a note",
				Action:  clicommands.GetContentAction(trilliumCli.conn),
			},
		},
	}

	logger.Info("Running: " + strings.Join(os.Args, " "))
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
	logger.Info("=== cli === trillium-cli === golang ===")
	logger.Info("=======================================")
	logger.Info("")

}
