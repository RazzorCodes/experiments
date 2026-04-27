package connector

import (
	"context"
	"io"
	"razzor/golang-mcp/internal/helpers"
	"razzor/golang-mcp/internal/ogen"
	"strings"
)

var triliumConv = helpers.NewTrilliumConverter()

func (conn *TrilliumConnector) getNoteById(ctx context.Context, noteId string) (*ogen.Note, error) {
	res, err := conn.client.GetNoteById(ctx, ogen.GetNoteByIdParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return nil, err
	}
	note, ok := res.(*ogen.Note)
	if !ok {
		return nil, ErrUnexpected
	}
	return note, nil
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
		noteType, _ := note.Type.Get()
		mime, _ := note.Mime.Get()

		resultMime := mime
		if mime == "text/html" {
			resultMime = "text/markdown"
		}

		snippets := helpers.ExtractWindowAroundKeywords(text, []string{keyword}, 128)
		hints := make([]Content, len(snippets))
		for i, s := range snippets {
			hints[i] = Content{Mime: resultMime, Value: s}
		}

		results = append(results, SearchResult{
			Title: title,
			Type:  string(noteType),
			Id:    id,
			Hints: hints,
		})
	}

	return results, nil
}

func (conn *TrilliumConnector) Content(noteId string) (*ContentResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()

	note, err := conn.getNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}

	body, err := conn.fetchContent(ctx, noteId)
	if err != nil {
		return nil, err
	}

	mime := note.Mime.Or("")
	contentValue := body
	contentMime := mime
	if mime == "text/html" {
		if md, err := triliumConv.ConvertString(body); err == nil {
			contentValue = md
			contentMime = "text/markdown"
		}
	}

	return &ContentResult{
		Title: note.Title.Or(""),
		Type:  string(note.Type.Or("")),
		Id:    string(note.NoteId.Or("")),
		Content: Content{
			Mime:  contentMime,
			Value: contentValue,
		},
	}, nil
}

func (conn *TrilliumConnector) Update(noteId string, content string) (*UpdateResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()

	note, err := conn.getNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}

	html, err := helpers.ConvertMDToHTML(content)
	if err != nil {
		return nil, err
	}

	err = conn.client.PutNoteContentById(
		ctx,
		ogen.PutNoteContentByIdReq{Data: strings.NewReader(html)},
		ogen.PutNoteContentByIdParams{NoteId: ogen.EntityId(noteId)},
	)
	if err != nil {
		return nil, err
	}

	return &UpdateResult{
		Title: note.Title.Or(""),
		Type:  string(note.Type.Or("")),
		Id:    string(note.NoteId.Or("")),
	}, nil
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

func (conn *TrilliumConnector) Delete(noteId string) (*DeleteResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()

	_, err := conn.client.DeleteNoteById(ctx, ogen.DeleteNoteByIdParams{
		NoteId: ogen.EntityId(noteId),
	})
	if err != nil {
		return nil, err
	}

	return &DeleteResult{Id: noteId}, nil
}
