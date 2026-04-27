package connector

import (
	"context"
	"errors"
	"razzor/golang-mcp/internal/config"
	"razzor/golang-mcp/internal/ogen"
	logger "razzor/golang-mcp/internal/utils"
)

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
