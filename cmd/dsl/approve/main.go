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

	base64Decoded, _ := base64.StdEncoding.DecodeString("CiQ3YTFkNGUyNy0xMzA5LTQ5OTAtYThmMC1jMTNiY2ZiMzkwNDASDGRzbC13b3JrZmxvdxokOTllZDY2MjktYmU5MS00YjJiLThmZmEtNTk3ZTBhZGZkMTdmIAUoATIBNUIQQXBwcm92YWxBY3Rpdml0eUoICAEQ4YJAGAE=")
	token := []byte(base64Decoded)
	state := "Approved"
	log.Debug().Msg("Trying to complete activity")
	err = workflowClient.CompleteActivity(context.Background(), token, state, nil)
	log.Debug().Str("state", state).Msg("Ccomplete activity")
	if err != nil {
		fmt.Printf("Failed to complete activity with error: %+v\n", err)
	} else {
		fmt.Printf("Successfully complete activity: %s\n", token)
	}
}
