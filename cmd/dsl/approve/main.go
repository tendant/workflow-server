package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
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
	state := "Approved"
	log.Debug().Msg("Trying to complete activity")
	// err = workflowClient.CompleteActivity(context.Background(), token, state, nil)
	namespace := "default"
	workflowId := "dsl-workflow"
	runId := "0b32fc81-2d78-4bec-beb1-f88b9d5d4c0d"
	activityId := "11"
	err = workflowClient.CompleteActivityByID(context.Background(), namespace, workflowId, runId, activityId, state, nil)
	log.Debug().Str("state", state).Msg("Ccomplete activity")
	if err != nil {
		fmt.Printf("Failed to complete activity with error: %+v\n", err)
	} else {
		fmt.Printf("Successfully complete activity: %s\n", token)
	}
}
