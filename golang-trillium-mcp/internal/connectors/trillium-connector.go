package connectors

import (
	"context"
	"errors"
	"io"
	"net/http"
	"razzor/trillium-mcp/internal/client"
	"razzor/trillium-mcp/internal/config"
	logger "razzor/trillium-mcp/internal/utils"
)

type TrilliumClient struct {
	Config config.Config
	Client *client.Client
}

var ErrHandshakeFailed = errors.New("Server handshake failed")
var ErrClientNotInit = errors.New("Client not initialized")

func NewTrilliumConnector(config config.Config) (*TrilliumClient, error) {
	newClient := &TrilliumClient{
		Config: config,
		Client: nil,
	}

	err := newClient.connect()
	if err != nil {
		return nil, err
	}

	err = newClient.test()
	if err != nil {
		return nil, err
	}

	return newClient, nil
}

func (conn *TrilliumClient) connect() error {
	newClient, err := client.NewClient(conn.Config.EtapiAddress)
	if err != nil {
		logger.Error("Could not connect to address: " + err.Error())
		return err
	}

	conn.Client = newClient

	return nil
}

func (conn *TrilliumClient) test() error {
	if conn.Client == nil {
		return ErrClientNotInit
	}
	ctx := context.Background()
	apikey := conn.Config.EtapiApikey
	address := conn.Config.EtapiAddress
	resp, err := conn.Client.GetAppInfo(
		ctx,
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+apikey)
			return nil
		})

	if err != nil {
		logger.Error("Could not establish connection to Trillium: " + address)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error("GetAppInfo: " + resp.Status + " " + string(body))
		return ErrHandshakeFailed
	}
	resp.Body.Close()

	return nil
}
