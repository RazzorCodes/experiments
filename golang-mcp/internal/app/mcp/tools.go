package mcpdefs

import (
	"context"
	"encoding/json"
	"razzor/golang-mcp/internal/config"
	connectors "razzor/golang-mcp/internal/connector"
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_notes",
		Description: "Search keywords in notes and note content",
	}, t.SearchKeyword)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_note_content",
		Description: "Get the full content of a note by its ID",
	}, t.GetContent)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_note",
		Description: "Update the content of a note by its ID. Content is Markdown and will be converted to HTML.",
	}, t.UpdateNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_note",
		Description: "Create a new note under a parent note. Content is Markdown and will be converted to HTML.",
	}, t.CreateNote)

	logger.Info("Sucessfully setup trillium mcp")

	t.Server = server

	return t, nil
}

// ---- param types ----

type SearchKeywordParams struct {
	Keyword string `json:"keyword"`
}

type NoteIdParams struct {
	NoteId string `json:"note_id"`
}

type UpdateNoteParams struct {
	NoteId  string `json:"note_id"`
	Content string `json:"content"`
}

type CreateNoteParams struct {
	ParentId string `json:"parent_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// ---- handlers ----

func (t TrilliumMcp) SearchKeyword(ctx context.Context, req *mcp.CallToolRequest, params SearchKeywordParams) (*mcp.CallToolResult, any, error) {
	results, err := t.conn.Search(params.Keyword)
	if err != nil {
		return nil, nil, err
	}

	var output string
	for _, r := range results {
		b, err := json.Marshal(r)
		if err != nil {
			return nil, nil, err
		}
		output += string(b) + "\n"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: output}},
	}, nil, nil
}

func (t TrilliumMcp) GetContent(ctx context.Context, req *mcp.CallToolRequest, params NoteIdParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Content(params.NoteId)
	if err != nil {
		return nil, nil, err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

func (t TrilliumMcp) UpdateNote(ctx context.Context, req *mcp.CallToolRequest, params UpdateNoteParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Update(params.NoteId, params.Content)
	if err != nil {
		return nil, nil, err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

func (t TrilliumMcp) CreateNote(ctx context.Context, req *mcp.CallToolRequest, params CreateNoteParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Create(params.ParentId, params.Title, params.Content)
	if err != nil {
		return nil, nil, err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}
