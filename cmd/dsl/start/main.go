package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tendant/workflow-server/dsl"
	"go.temporal.io/sdk/client"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create Temporal client")
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "dsl-workflow",
		TaskQueue: dsl.DSLWorkflowTaskQueue,
	}
	log.Info().Msg("workflow options:")

	// Start the Workflow
	args := dsl.DSLWorkflowArgs{
		ExpenseId: "first",
	}
	log.Info().Msg("exeucte workflow")
	we, err := c.ExecuteWorkflow(context.Background(), options, dsl.DSLWorkflow, args)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to complete Workflow")
	}
	log.Info().Str("workflowID", we.GetID()).Str("runID", we.GetRunID()).Msg("Started Workflow")

	// Get the results
	var wf string
	log.Info().Msg("Getting workflow result...")
	err = we.Get(context.Background(), &wf)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get Workflow result")
	}

	printResults(wf, we.GetID(), we.GetRunID())
}

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", result)
}
