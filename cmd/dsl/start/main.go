package main

import (
	"context"

	"github.com/tendant/workflow-server/dsl"
	"go.temporal.io/sdk/client"
	"golang.org/x/exp/slog"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("unable to create Temporal client", "err", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "dsl-workflow",
		TaskQueue: dsl.DSLWorkflowTaskQueue,
	}
	slog.Info("workflow options:")

	// Start the Workflow
	args := dsl.DSLWorkflowArgs{
		ExpenseId: "first",
	}
	slog.Info("exeucte workflow")
	we, err := c.ExecuteWorkflow(context.Background(), options, dsl.DSLWorkflow, args)
	if err != nil {
		slog.Error("unable to complete Workflow", "err", err)
	}
	slog.Info("Started Workflow", "workflowID", we.GetID(), "runID", we.GetRunID())

	// Get the results
	var wf string
	slog.Info("Getting workflow result...")
	err = we.Get(context.Background(), &wf)
	if err != nil {
		slog.Error("unable to get Workflow result", "err", err)
	}

	printResults(wf, we.GetID(), we.GetRunID())
}

func printResults(result string, workflowID, runID string) {
	slog.Info("Workflow Results:", "result", result, "workflowID", workflowID, "runID", runID)
}
