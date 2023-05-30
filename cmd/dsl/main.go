package main

import (
	"github.com/rs/zerolog/log"
	"github.com/tendant/workflow-server/dsl"
)

func main() {

	log.Debug().Msg("Starting...")
	workflow, err := dsl.ParseWorkflow("./expense.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed parsing workflow")
	}
	log.Info().Msg("Workflow info:")
	log.Info().Str("Start state", workflow.Start.StateName).Send()

}
