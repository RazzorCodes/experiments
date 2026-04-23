package mcpdefs

import (
	"context"
	"encoding/json"
	"razzor/golang-mcp/internal/config"
	connectors "razzor/golang-mcp/internal/connectors/trillium"
	logger "razzor/golang-mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TrilliumMcp struct {
	Server *mcp.Server
	conn   *connectors.TrilliumConnector
}

func (t TrilliumMcp) NewTrilliumMcp() (TrilliumMcp, error) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "trillium-notes",
		Version: "0.0.1",
	}, nil)

	t.Server = server

	logger.Info("Setting up...")

	appconf, err := config.GetConfig()
	if err != nil {
		logger.Info("Config was not loaded sucessfully")
	}

	logger.Info("AppVersion: " + appconf.AppVersion)
	logger.Info("EtapiAddress: " + appconf.EtapiAddress)

	t.conn, err = connectors.NewTrilliumConnector(appconf)
	if err != nil {
		logger.Fatal("")
	}

	logger.Info("Adding tools")

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "search_notes",
			Description: "Search keywords in notes and note content",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type":        "string",
						"description": "Search keyword",
					},
				},
				"required": []string{"keyword"},
			},
		},
		t.SearchKeyword)

	logger.Info("Sucessfully setup trillium mcp")

	t.Server = server

	return t, nil
}

type SearchKeywordParams struct {
	Keyword string `json:"keyword"`
}

func (t TrilliumMcp) SearchKeyword(ctx context.Context, req *mcp.CallToolRequest, params SearchKeywordParams) (*mcp.CallToolResult, any, error) {
	result, _ := t.conn.Search(params.Keyword)

	var output string
	for i := range result {
		jsonBytes, err := json.Marshal(result[i])
		if err != nil {
			return nil, nil, err
		}
		output += string(jsonBytes) + "\n"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: output}},
	}, nil, nil
}
