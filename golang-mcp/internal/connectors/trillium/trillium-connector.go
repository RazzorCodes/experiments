package connectors

import (
	"context"
	"errors"
	"net/http"
	"razzor/golang-mcp/internal/client"
	"razzor/golang-mcp/internal/config"
	"razzor/golang-mcp/internal/helpers"
	logger "razzor/golang-mcp/internal/utils"
	"strings"

)

var triliumConv = newTrilliumConverter()

type TrilliumConnector struct {
	Client *client.ClientWithResponses

	requestEditor func(ctx context.Context, req *http.Request) error
}

var ErrHandshakeFailed = errors.New("Server handshake failed")
var ErrClientNotInit = errors.New("Client not initialized")

func NewTrilliumConnector(config config.Config) (*TrilliumConnector, error) {
	newConnector := &TrilliumConnector{
		Client: nil,
		requestEditor: func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+config.EtapiApikey)
			req.Header.Set("Accept", "text/html")
			return nil
		},
	}

	err := newConnector.connect(config.EtapiAddress)
	if err != nil {
		return nil, err
	}

	err = newConnector.test()
	if err != nil {
		return nil, err
	}

	return newConnector, nil
}

func (conn *TrilliumConnector) connect(address string) error {
	newClient, err := client.NewClientWithResponses(address)
	if err != nil {
		logger.Error("Could not connect to address: " + err.Error())
		return err
	}

	conn.Client = newClient

	return nil
}

func (conn *TrilliumConnector) test() error {
	if conn.Client == nil {
		return ErrClientNotInit
	}

	ctx := context.Background()
	resp, err := conn.Client.GetAppInfoWithResponse(ctx, conn.requestEditor)
	if err != nil {
		logger.Error("Could not establish connection to Trillium")
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error("GetAppInfo: " + resp.Status() + " " + string(resp.Body))
		return ErrHandshakeFailed
	}

	return nil
}

func firstNLines(text string, n int) string {
	lines := strings.Split(text, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}

func truncate(text string, maxChars int) string {
	if maxChars > 0 && len(text) > maxChars {
		return text[:maxChars]
	}
	return text
}

func (conn *TrilliumConnector) Search(keyword string) ([]SearchResult, error) {
	if conn == nil || conn.Client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()
	resp, err := conn.Client.SearchNotesWithResponse(
		ctx,
		&client.SearchNotesParams{Search: keyword, FastSearch: func() *bool { b := false; return &b }()},
		conn.requestEditor)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error("SearchNotes failed: " + resp.Status())
		return nil, ErrHandshakeFailed
	}

	if resp.ApplicationjsonCharsetUtf8200 != nil {
		var results []SearchResult
		for _, note := range resp.ApplicationjsonCharsetUtf8200.Results {
			content, err :=
				conn.Client.GetNoteContentWithResponse(
					ctx,
					*note.NoteId,
					conn.requestEditor)

			if err != nil || content.StatusCode() != http.StatusOK {
				continue
			}

			var text string
			if *note.Mime == "text/html" {
				text, _ = triliumConv.ConvertString(string(content.Body))
			} else {
				text = string(content.Body)
			}

			noteData := SearchResult{
				Title: *note.Title,
				Type:  *note.Mime,
				Id:    *note.NoteId,
				Hints: helpers.ExtractWindowAroundKeywords(
					string(text),
					[]string{keyword},
					128,
				),
			}

			results = append(results, noteData)
		}

		return results, nil
	}

	return nil, nil
}

func (conn *TrilliumConnector) Content(noteId *client.EntityId) (string, error) {
	if conn == nil || noteId == nil {
		return "", ErrClientNotInit
	}

	ctx := context.Background()
	note, err := conn.Client.GetNoteByIdWithResponse(ctx, *noteId, conn.requestEditor)
	if err != nil {
		return "", err
	}
	if note.StatusCode() != http.StatusOK {
		return "", errors.New(
			note.ApplicationjsonCharsetUtf8Default.Code + " : " +
				note.ApplicationjsonCharsetUtf8Default.Message)
	}

	content, err :=
		conn.Client.GetNoteContentWithResponse(
			ctx,
			*noteId,
			conn.requestEditor)
	if err != nil {
		return "", err
	}

	mime := ""
	if note.ApplicationjsonCharsetUtf8200 != nil && note.ApplicationjsonCharsetUtf8200.Mime != nil {
		mime = *note.ApplicationjsonCharsetUtf8200.Mime
	}

	body := string(content.Body)
	if mime == "text/html" {
		if md, err := triliumConv.ConvertString(body); err == nil {
			return md, nil
		}
	}
	return body, nil
}
