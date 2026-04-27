package mcpdefs

import (
	"context"
	"encoding/json"
	"razzor/golang-mcp/internal/config"
	connector "razzor/golang-mcp/internal/connector"
	logger "razzor/golang-mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TrilliumMcp struct {
	Server *mcp.Server
	conn   *connector.TrilliumConnector
}

func NewTrilliumMcp() (TrilliumMcp, error) {
	var t TrilliumMcp

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "trillium-notes",
		Version: "0.0.1",
	}, nil)

	logger.Info("Setting up...")

	appconf, err := config.GetConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("AppVersion: " + appconf.AppVersion)
	logger.Info("EtapiAddress: " + appconf.EtapiAddress)

	t.conn, err = connector.NewTrilliumConnector(appconf)
	if err != nil {
		logger.Fatal(err.Error())
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "move_note",
		Description: "Move a note to a new parent. Removes all existing parent branches and places the note under the new parent.",
	}, t.MoveNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_note",
		Description: "Delete a note and all its children by note ID.",
	}, t.DeleteNote)

	logger.Info("Sucessfully setup trillium mcp")

	t.Server = server
	return t, nil
}

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

type MoveNoteParams struct {
	NoteId      string `json:"note_id"`
	NewParentId string `json:"new_parent_id"`
}

func textResult(v any) (*mcp.CallToolResult, any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

func (t TrilliumMcp) SearchKeyword(ctx context.Context, req *mcp.CallToolRequest, params SearchKeywordParams) (*mcp.CallToolResult, any, error) {
	results, err := t.conn.Search(params.Keyword)
	if err != nil {
		return nil, nil, err
	}
	return textResult(results)
}

func (t TrilliumMcp) GetContent(ctx context.Context, req *mcp.CallToolRequest, params NoteIdParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Content(params.NoteId)
	if err != nil {
		return nil, nil, err
	}
	return textResult(res)
}

func (t TrilliumMcp) UpdateNote(ctx context.Context, req *mcp.CallToolRequest, params UpdateNoteParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Update(params.NoteId, params.Content)
	if err != nil {
		return nil, nil, err
	}
	return textResult(res)
}

func (t TrilliumMcp) CreateNote(ctx context.Context, req *mcp.CallToolRequest, params CreateNoteParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Create(params.ParentId, params.Title, params.Content)
	if err != nil {
		return nil, nil, err
	}
	return textResult(res)
}

func (t TrilliumMcp) MoveNote(ctx context.Context, req *mcp.CallToolRequest, params MoveNoteParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Move(params.NoteId, params.NewParentId)
	if err != nil {
		return nil, nil, err
	}
	return textResult(res)
}

func (t TrilliumMcp) DeleteNote(ctx context.Context, req *mcp.CallToolRequest, params NoteIdParams) (*mcp.CallToolResult, any, error) {
	res, err := t.conn.Delete(params.NoteId)
	if err != nil {
		return nil, nil, err
	}
	return textResult(res)
}
