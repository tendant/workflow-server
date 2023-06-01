package main

import (
	"github.com/rs/zerolog/log"
	"github.com/tendant/workflow-server/dsl"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	log.Debug().Msg("Starting worker...")
	// workflow, err := dsl.ParseWorkflow("./expense.yaml")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed parsing workflow")
	// }
	// log.Info().Msg("Workflow info:")

	// log.Info().Str("Start state", workflow.Start.StateName).Send()

	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatal().AnErr("err", err).Msg("unable to create Temporal client for woker")
	}
	defer c.Close()

	w := worker.New(c, dsl.DSLWorkflowTaskQueue, worker.Options{})
	w.RegisterWorkflow(dsl.DSLWorkflow)
	w.RegisterActivity(dsl.ApprovalActivity)

	// start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal().AnErr("err", err).Msg("Unable to start Worker")
	}

}
