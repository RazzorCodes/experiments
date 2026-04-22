package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	client "razzor/trillium-mcp/internal/client"
	logger "razzor/trillium-mcp/internal/utils"
	"strconv"
)

const (
	UserAgentName    = "trillium-mcp"
	UserAgentVersion = "0.0.1"
	UserAgent        = UserAgentName + "/" + UserAgentVersion
	// ---
	TrilliumAddress  = "http://notes.lan/etapi"
	TrilliumApiToken = "CBkpReV6TEM6_FumFTgWmj1N3BWVkRnYd9RKtPJVdhkGeRXpInKnXfns="
)

func testTrilliumConnection() bool {

	trilliumClient, _ := client.NewClient(TrilliumAddress)
	ctx := context.Background()

	resp, err := trilliumClient.GetAppInfo(ctx, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+TrilliumApiToken)
		return nil
	})
	if err != nil {
		logger.Get().Fatal("Could not establish connection to Trillium: " + TrilliumAddress)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Get().Fatal("GetAppInfo: " + resp.Status + " " + string(body))
	}
	defer resp.Body.Close()

	var appInfo client.AppInfo
	json.NewDecoder(resp.Body).Decode(&appInfo)

	resp, err = trilliumClient.SearchNotes(ctx, &client.SearchNotesParams{Search: "HomeReef"},
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+TrilliumApiToken)
			return nil
		},
	)

	if err != nil {
		logger.Get().Fatal("Could not retreive data")
	}

	body, _ := io.ReadAll(resp.Body)

	logger.Get().Info(strconv.Itoa(resp.StatusCode) + string(body))

	return true
}

func main() {
	logger.Get().Info("=== mcp === trillium-mcp === golang ===")
	logger.Get().Info("mcp: " + UserAgent)
	logger.Get().Info("Initializing...")

	if testTrilliumConnection() {
		logger.Get().Info("Ready!")
	} else {
		logger.Get().Fatal("Crash")
	}
}
