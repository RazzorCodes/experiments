package clicommands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	connectors "razzor/golang-mcp/internal/connectors/trillium"

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

		id := cmd.Args().First()
		res, err := conn.Content(&id)
		if err != nil {
			return err
		}

		fmt.Println(res)
		return nil
	}
}

// const (
// 	search = clisdk.Command{
// 	{
//         Name: "search",
//         Flags: []cli.Flag{
//             &cli.StringFlag{
//                 Name:  "format",
//                 Value: "raw",
//             },
//         },
//         Action: func(c *cli.Context) error {
//             modifier := c.String("modifier")
//             quiet := c.Bool("q")
//
//             log.Printf("modifier=%s quiet=%v\n", modifier, quiet)
//             return nil
//         },
//     },
// )
