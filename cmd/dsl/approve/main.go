package main

import (
	"context"
	"encoding/base64"

	"go.temporal.io/sdk/client"
	"golang.org/x/exp/slog"
)

func main() {

	// The client is a heavyweight object that should be created once per process.
	workflowClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		panic(err)
	}

	base64Decoded, _ := base64.StdEncoding.DecodeString("CiQ3YTFkNGUyNy0xMzA5LTQ5OTAtYThmMC1jMTNiY2ZiMzkwNDASDGRzbC13b3JrZmxvdxokOTllZDY2MjktYmU5MS00YjJiLThmZmEtNTk3ZTBhZGZkMTdmIAsoATICMTFCEEFwcHJvdmFsQWN0aXZpdHlKCAgBEPGCQBgB")
	token := []byte(base64Decoded)
	// state := "Approved"
	state := "Declined"
	slog.Debug("Trying to complete activity")
	// err = workflowClient.CompleteActivity(context.Background(), token, state, nil)
	namespace := "default"
	workflowId := "tx-approval-1"
	runId := "bda7cf4b-0d62-4736-bbd1-a94233019fb4"
	activityId := "11"
	err = workflowClient.CompleteActivityByID(context.Background(), namespace, workflowId, runId, activityId, state, nil)
	slog.Debug("Ccomplete activity", "state", state)
	if err != nil {
		slog.Error("Failed to complete activity with error: %+v\n", "err", err)
	} else {
		slog.Info("Successfully complete activity", "token", token)
	}
}
