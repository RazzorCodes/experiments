package connectors

import (
	"context"
	"errors"
	"net/http"
	"razzor/golang-mcp/internal/client"
	"razzor/golang-mcp/internal/config"
	logger "razzor/golang-mcp/internal/utils"
	"strings"

	"github.com/jaytaylor/html2text"
)

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
	if len(text) > maxChars {
		return text[:maxChars]
	}
	return text
}

func (conn *TrilliumConnector) Search(keyword string) ([]NoteResult, error) {
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
		var result []NoteResult
		for _, note := range resp.ApplicationjsonCharsetUtf8200.Results {
			content, err := conn.Client.GetNoteContentWithResponse(ctx, *note.NoteId, conn.requestEditor)
			if err != nil || content.StatusCode() != http.StatusOK {
				continue
			}

			var text string
			if *note.Mime == "text/html" {
				text, _ = html2text.FromString(string(content.Body))
			} else {
				text = string(content.Body)
			}

			noteData := NoteResult{
				Title:    *note.Title,
				Type:     *note.Mime,
				Id:       *note.BlobId,
				Contents: truncate(text, 1024),
			}

			result = append(result, noteData)
		}

		return result, nil
	}

	return nil, nil
}
