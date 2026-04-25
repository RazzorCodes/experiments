package connector

import (
	"context"
	"razzor/golang-mcp/internal/ogen"

	"github.com/ogen-go/ogen/ogenerrors"
)

type etapiAuth struct {
	token string
}

func (a *etapiAuth) EtapiTokenAuth(_ context.Context, _ ogen.OperationName) (ogen.EtapiTokenAuth, error) {
	return ogen.EtapiTokenAuth{APIKey: "Bearer " + a.token}, nil
}

func (a *etapiAuth) EtapiBasicAuth(_ context.Context, _ ogen.OperationName) (ogen.EtapiBasicAuth, error) {
	return ogen.EtapiBasicAuth{}, ogenerrors.ErrSkipClientSecurity
}
