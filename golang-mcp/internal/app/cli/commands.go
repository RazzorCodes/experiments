package clicommands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	connector "razzor/golang-mcp/internal/connector"

	clisdk "github.com/urfave/cli/v3"
)

var ErrInvalidParam = errors.New("Parameter value invalid")

func contentFromFlags(cmd *clisdk.Command) (string, error) {
	if file := cmd.String("path"); file != "" {
		data, err := os.ReadFile(file)
		return string(data), err
	}
	return cmd.String("content"), nil
}

func printJSON(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func GetSearchAction(conn *connector.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		res, err := conn.Search(cmd.Args().First())
		if err != nil {
			return err
		}

		for _, r := range res {
			if err := printJSON(r); err != nil {
				return err
			}
		}
		return nil
	}
}

func GetContentAction(conn *connector.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		res, err := conn.Content(cmd.String("id"))
		if err != nil {
			return err
		}

		return printJSON(res)
	}
}

func GetUpdateAction(conn *connector.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		content, err := contentFromFlags(cmd)
		if err != nil {
			return err
		}

		res, err := conn.Update(cmd.String("id"), content)
		if err != nil {
			return err
		}

		return printJSON(res)
	}
}

func GetMoveAction(conn *connector.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		res, err := conn.Move(cmd.String("id"), cmd.String("parent"))
		if err != nil {
			return err
		}

		return printJSON(res)
	}
}

func GetDeleteAction(conn *connector.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		res, err := conn.Delete(cmd.String("id"))
		if err != nil {
			return err
		}

		return printJSON(res)
	}
}

func GetCreateAction(conn *connector.TrilliumConnector) clisdk.ActionFunc {
	return func(ctx context.Context, cmd *clisdk.Command) error {
		if conn == nil {
			return ErrInvalidParam
		}

		content, err := contentFromFlags(cmd)
		if err != nil {
			return err
		}

		res, err := conn.Create(cmd.String("parent"), cmd.String("title"), content)
		if err != nil {
			return err
		}

		return printJSON(res)
	}
}
