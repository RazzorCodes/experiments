package connectors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"razzor/golang-mcp/internal/config"
	"razzor/golang-mcp/internal/helpers"
	"razzor/golang-mcp/internal/ogen"
	logger "razzor/golang-mcp/internal/utils"
)

var triliumConv = helpers.NewTrilliumConverter()

var ErrHandshakeFailed = errors.New("server handshake failed")
var ErrClientNotInit = errors.New("client not initialized")

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
		return nil, fmt.Errorf("unexpected search response type: %T", res)
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

func (conn *TrilliumConnector) Content(noteId string) (string, error) {
	if conn == nil || conn.client == nil {
		return "", ErrClientNotInit
	}

	ctx := context.Background()

	noteRes, err := conn.client.GetNoteById(ctx, ogen.GetNoteByIdParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return "", err
	}
	note, ok := noteRes.(*ogen.Note)
	if !ok {
		return "", fmt.Errorf("unexpected note response type: %T", noteRes)
	}

	body, err := conn.fetchContent(ctx, noteId)
	if err != nil {
		return "", err
	}

	if mime, ok := note.Mime.Get(); ok && mime == "text/html" {
		if md, err := triliumConv.ConvertString(body); err == nil {
			return md, nil
		}
	}
	return body, nil
}

func (conn *TrilliumConnector) fetchContent(ctx context.Context, noteId string) (string, error) {
	res, err := conn.client.GetNoteContent(ctx, ogen.GetNoteContentParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(res.Data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// etapiAuth implements ogen.SecuritySource using the ETAPI token.
type etapiAuth struct {
	token string
}

func (a *etapiAuth) EtapiTokenAuth(_ context.Context, _ ogen.OperationName) (ogen.EtapiTokenAuth, error) {
	return ogen.EtapiTokenAuth{APIKey: "Bearer " + a.token}, nil
}

func (a *etapiAuth) EtapiBasicAuth(_ context.Context, _ ogen.OperationName) (ogen.EtapiBasicAuth, error) {
	return ogen.EtapiBasicAuth{}, nil
}
