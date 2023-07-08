package main

import (
	"context"
	"fmt"
	"os"

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

	filename := "./samples/bankingtransactions.yaml"
	dslStr, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Failed reading workflow DSL from file.", "file", filename)
		os.Exit(-1)
	}
	wfType := "tx"
	entityType := "bankingtransaction"
	entityId := "1"
	id := fmt.Sprintf("%s-%s-%s", wfType, entityType, entityId)
	args := dsl.DSLWorkflowArgs{
		Id:         id,
		Type:       wfType,
		EntityType: entityType,
		EntityId:   entityId,
		DSLStr:     string(dslStr),
	}

	slog.Info("exeucte workflow")
	// Start the Workflow
	options := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: dsl.DSLWorkflowTaskQueue,
	}
	slog.Info("workflow options:", "options", options)

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
