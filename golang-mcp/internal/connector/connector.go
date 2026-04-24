package connectors

import (
	"context"
	"errors"
	"io"
	"razzor/golang-mcp/internal/config"
	"razzor/golang-mcp/internal/helpers"
	"razzor/golang-mcp/internal/ogen"
	logger "razzor/golang-mcp/internal/utils"
	"strings"

	"github.com/ogen-go/ogen/ogenerrors"
)

var triliumConv = helpers.NewTrilliumConverter()

var ErrHandshakeFailed = errors.New("server handshake failed")
var ErrClientNotInit = errors.New("client not initialized")
var ErrUnexpected = errors.New("unexpected value")

type TrilliumConnector struct {
	client *ogen.Client
}

func NewTrilliumConnector(cfg config.Config) (*TrilliumConnector, error) {
	c, err := ogen.NewClient(cfg.EtapiAddress, &etapiAuth{token: cfg.EtapiApikey})
	if err != nil {
		return nil, err
	}

	conn := &TrilliumConnector{client: c}

	if err := conn.test(); err != nil {
		return nil, err
	}

	return conn, nil
}

func (conn *TrilliumConnector) test() error {
	res, err := conn.client.GetAppInfo(context.Background())
	if err != nil {
		logger.Error("Could not connect to Trillium: " + err.Error())
		return err
	}
	if _, ok := res.(*ogen.AppInfo); !ok {
		return ErrHandshakeFailed
	}
	return nil
}

func (conn *TrilliumConnector) Search(keyword string) ([]SearchResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()
	res, err := conn.client.SearchNotes(ctx, ogen.SearchNotesParams{
		Search: keyword,
	})
	if err != nil {
		return nil, err
	}

	sr, ok := res.(*ogen.SearchResponse)
	if !ok {
		return nil, ErrUnexpected
	}

	var results []SearchResult
	for _, note := range sr.Results {
		noteId, ok := note.NoteId.Get()
		if !ok {
			continue
		}
		id := string(noteId)

		body, err := conn.fetchContent(ctx, id)
		if err != nil {
			continue
		}

		text := body
		if mime, ok := note.Mime.Get(); ok && mime == "text/html" {
			if md, err := triliumConv.ConvertString(body); err == nil {
				text = md
			}
		}

		title, _ := note.Title.Get()
		mime, _ := note.Mime.Get()

		results = append(results, SearchResult{
			Title: title,
			Type:  mime,
			Id:    id,
			Hints: helpers.ExtractWindowAroundKeywords(text, []string{keyword}, 128),
		})
	}

	return results, nil
}

func (conn *TrilliumConnector) Content(noteId string) (*ContentResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()

	noteRes, err := conn.client.GetNoteById(ctx, ogen.GetNoteByIdParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return nil, err
	}
	note, ok := noteRes.(*ogen.Note)
	if !ok {
		return nil, ErrUnexpected
	}

	body, err := conn.fetchContent(ctx, noteId)
	if err != nil {
		return nil, err
	}

	result := ContentResult{
		Title:   note.Title.Or(""),
		Type:    note.Mime.Or(""),
		Id:      string(note.NoteId.Or("")),
		Content: body,
	}

	if result.Type == "text/html" {
		if md, err := triliumConv.ConvertString(body); err == nil {
			result.Content = md
		}
	}

	return &result, nil
}

func (conn *TrilliumConnector) fetchContent(ctx context.Context, noteId string) (string, error) {
	res, err := conn.client.GetNoteContent(ctx, ogen.GetNoteContentParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(res.Response.Data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (conn *TrilliumConnector) Update(noteId string, content string) (*UpdateResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()

	noteRes, err := conn.client.GetNoteById(ctx, ogen.GetNoteByIdParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return nil, err
	}
	note, ok := noteRes.(*ogen.Note)
	if !ok {
		return nil, ErrUnexpected
	}

	html, err := helpers.ConvertMDToHTML(content)
	if err != nil {
		return nil, err
	}

	success, err := conn.updateContent(ctx, noteId, html)
	if err != nil {
		return nil, err
	}

	result := UpdateResult{
		Title:   note.Title.Or(""),
		Type:    note.Mime.Or(""),
		Id:      string(note.NoteId.Or("")),
		Success: success,
	}

	return &result, nil
}

func (conn *TrilliumConnector) updateContent(ctx context.Context, noteId string, content string) (bool, error) {
	err := conn.client.PutNoteContentById(
		ctx,
		ogen.PutNoteContentByIdReq{Data: strings.NewReader(content)},
		ogen.PutNoteContentByIdParams{NoteId: ogen.EntityId(noteId)})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (conn *TrilliumConnector) Create(parentId string, title string, content string) (*CreateResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	html, err := helpers.ConvertMDToHTML(content)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	res, err := conn.client.CreateNote(ctx, &ogen.CreateNoteDef{
		ParentNoteId: ogen.EntityId(parentId),
		Title:        title,
		Type:         ogen.CreateNoteDefTypeText,
		Mime:         ogen.NewOptString("text/html"),
		Content:      html,
	})
	if err != nil {
		return nil, err
	}

	nwb, ok := res.(*ogen.NoteWithBranch)
	if !ok {
		return nil, ErrUnexpected
	}

	note, _ := nwb.GetNote().Get()
	return &CreateResult{
		Title: note.Title.Or(""),
		Id:    string(note.NoteId.Or("")),
	}, nil
}

type etapiAuth struct {
	token string
}

func (a *etapiAuth) EtapiTokenAuth(_ context.Context, _ ogen.OperationName) (ogen.EtapiTokenAuth, error) {
	return ogen.EtapiTokenAuth{APIKey: "Bearer " + a.token}, nil
}

func (a *etapiAuth) EtapiBasicAuth(_ context.Context, _ ogen.OperationName) (ogen.EtapiBasicAuth, error) {
	return ogen.EtapiBasicAuth{}, ogenerrors.ErrSkipClientSecurity
}
