package clicommands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	connectors "razzor/golang-mcp/internal/connector"

	clisdk "github.com/urfave/cli/v3"
)

var ErrInvalidParam = errors.New("Parameter value invalid")

func GetSearchAction(conn *connectors.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		res, err := conn.Search(cmd.Args().First())
		if err != nil {
			return err
		}

		for _, r := range res {
			b, _ := json.Marshal(r)
			fmt.Println(string(b))
		}
		return nil
	}
}

func GetContentAction(conn *connectors.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		noteId := cmd.String("id")

		res, err := conn.Content(noteId)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(res)
		fmt.Println(string(b))
		return nil
	}
}

func GetUpdateAction(conn *connectors.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		noteId := cmd.String("id")
		file := cmd.String("path")
		content := ""
		if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			content = string(data)
		} else {
			content = cmd.String("content")
		}

		res, err := conn.Update(noteId, content)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(res)
		fmt.Println(string(b))
		return nil
	}
}

func GetCreateAction(conn *connectors.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		title := cmd.String("title")
		parent := cmd.String("parent")
		file := cmd.String("path")
		content := ""
		if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			content = string(data)
		} else {
			content = cmd.String("content")
		}

		res, err := conn.Create(parent, title, content)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(res)
		fmt.Println(string(b))
		return nil
	}
}
