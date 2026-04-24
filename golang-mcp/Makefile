OGEN      := $(shell go env GOPATH)/bin/ogen
SPEC      := api/trillium-openapi.yaml
OGEN_CFG  := api/ogen.yml
OGEN_OUT  := internal/ogen

.PHONY: init generate build clean

init:
	go mod download
	go mod tidy

generate:
	$(OGEN) --config $(OGEN_CFG) --target $(OGEN_OUT) --package ogen $(SPEC)

build:
	go build -o bin/ ./cmd/...

clean:
	rm -rf $(OGEN_OUT)
