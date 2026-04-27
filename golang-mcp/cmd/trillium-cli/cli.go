package main

import (
	"context"
	"log"
	"os"
	clicommands "razzor/golang-mcp/internal/app/cli"
	"razzor/golang-mcp/internal/config"
	connector "razzor/golang-mcp/internal/connector"
	logger "razzor/golang-mcp/internal/utils"
	"strings"

	"github.com/urfave/cli/v3" // imports as package "cli"
)

type TrilliumCli struct {
	conn *connector.TrilliumConnector
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
	t.conn, err = connector.NewTrilliumConnector(appconf)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return t, nil
}

func main() {
	logger.Info("=======================================")
	logger.Info("=== cli === trillium-cli === golang ===")

	trilliumCli, err := NewTrilliumCli()
	if err != nil {
		logger.Fatal(err.Error())
	}

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "search",
				Usage:  "search for a keyword",
				Action: clicommands.GetSearchAction(trilliumCli.conn),
			},
			{
				Name: "content",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
				},
				Usage:  "get content of a note",
				Action: clicommands.GetContentAction(trilliumCli.conn),
			},
			{
				Name: "update",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
					&cli.StringFlag{
						Name: "path",
					},
					&cli.StringFlag{
						Name: "content",
					},
				},
				Usage:  "update content of a note",
				Action: clicommands.GetUpdateAction(trilliumCli.conn),
			},
			{
				Name: "move",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
					&cli.StringFlag{
						Name: "parent",
					},
				},
				Usage:  "move a note to a new parent",
				Action: clicommands.GetMoveAction(trilliumCli.conn),
			},
			{
				Name: "delete",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
				},
				Usage:  "delete a note by id",
				Action: clicommands.GetDeleteAction(trilliumCli.conn),
			},
			{
				Name: "add",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "title",
					},
					&cli.StringFlag{
						Name: "parent",
					},
					&cli.StringFlag{
						Name: "path",
					},
					&cli.StringFlag{
						Name: "content",
					},
				},
				Usage:  "add a new note",
				Action: clicommands.GetCreateAction(trilliumCli.conn),
			},
		},
	}

	logger.Info("Running: " + strings.Join(os.Args, " "))
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Error(err.Error())
		log.Fatal(err)
	}
	logger.Info("=== cli === trillium-cli === golang ===")
	logger.Info("=======================================")
	logger.Info("")

}
