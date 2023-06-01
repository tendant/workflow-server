package main

import (
	"context"
	"fmt"

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

	token := []byte("")
	state := "Approved"
	err = workflowClient.CompleteActivity(context.Background(), token, state, nil)
	if err != nil {
		fmt.Printf("Failed to complete activity with error: %+v\n", err)
	} else {
		fmt.Printf("Successfully complete activity: %s\n", token)
	}
}
