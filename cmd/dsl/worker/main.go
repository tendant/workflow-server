package main

import (
	"github.com/tendant/workflow-server/dsl"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"golang.org/x/exp/slog"
)

func main() {

	slog.Debug("Starting worker...")
	// workflow, err := dsl.ParseWorkflow("./expense.yaml")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed parsing workflow")
	// }
	// log.Info().Msg("Workflow info:")

	// log.Info().Str("Start state", workflow.Start.StateName).Send()

	c, err := client.Dial(client.Options{})

	if err != nil {
		slog.Error("unable to create Temporal client for woker", "err", err)
	}
	defer c.Close()

	w := worker.New(c, dsl.DSLWorkflowTaskQueue, worker.Options{})
	w.RegisterWorkflow(dsl.DSLWorkflow)
	regActivityOptions := activity.RegisterOptions{
		Name: "TransactionApprovalActivity",
	}
	w.RegisterActivityWithOptions(dsl.TransactionApprovalActivity, regActivityOptions)

	// start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("Unable to start Worker", "err", err)
	}

}
