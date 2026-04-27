package connector

import (
	"context"
	"razzor/golang-mcp/internal/ogen"
)

func (conn *TrilliumConnector) Move(noteId, newParentId string) (*MoveResult, error) {
	if conn == nil || conn.client == nil {
		return nil, ErrClientNotInit
	}

	ctx := context.Background()

	note, err := conn.getNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}

	res, err := conn.client.PostBranch(ctx, &ogen.Branch{
		NoteId:       ogen.NewOptEntityId(ogen.EntityId(noteId)),
		ParentNoteId: ogen.NewOptEntityId(ogen.EntityId(newParentId)),
	})
	if err != nil {
		return nil, err
	}

	newBranch, ok := res.(*ogen.PostBranchOK)
	if !ok {
		return nil, ErrUnexpected
	}

	for _, branchId := range note.ParentBranchIds {
		conn.client.DeleteBranchById(ctx, ogen.DeleteBranchByIdParams{
			BranchId: ogen.EntityId(branchId),
		})
	}

	return &MoveResult{
		Title:    note.Title.Or(""),
		Type:     string(note.Type.Or("")),
		Id:       noteId,
		BranchId: string(newBranch.BranchId.Or("")),
	}, nil
}
