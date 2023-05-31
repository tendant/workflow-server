package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tendant/workflow-server/dsl"
	"go.temporal.io/sdk/client"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "dsl-workflow",
		TaskQueue: dsl.DSLWorkflowTaskQueue,
	}

	// Start the Workflow
	args := &dsl.DSLWorkflowArgs{
		ExpenseId: "first",
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, dsl.DSLWorkflow, args)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var greeting string
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(greeting, we.GetID(), we.GetRunID())
}

func printResults(greeting string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", greeting)
}
