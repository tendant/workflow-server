package dsl

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

type WorkflowClient struct {
	wfc client.Client
	log zerolog.Logger
}

type WorkflowRunAct struct {
	Namespace  string
	WorkflowId string
	RunId      string
	ActivityId string
	State      string
	Err        error
}

func NewWorkflowClient() (*WorkflowClient, error) {
	// The client is a heavyweight object that should be created once per process.
	workflowClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	log := log.Logger
	if err != nil {
		log.Err(err).Msg("Failed to dial Temporal server")
		return nil, err
	}
	return &WorkflowClient{workflowClient, log}, nil
}

func (w *WorkflowClient) CompleteActivityByID(runact WorkflowRunAct) error {
	w.log.Debug().Msg("Trying to complete activity")
	namespace := runact.Namespace   // "default"
	workflowId := runact.WorkflowId // "dsl-workflow"
	runId := runact.RunId           // "0b32fc81-2d78-4bec-beb1-f88b9d5d4c0d"
	activityId := runact.ActivityId // "11"
	state := runact.State           // "Approved"
	err := runact.Err
	err = w.wfc.CompleteActivityByID(context.Background(), namespace, workflowId, runId, activityId, state, err)
	w.log.Debug().Str("state", state).Msg("Ccomplete activity")
	if err != nil {
		w.log.Error().AnErr("Failed to complete activity with error", err)
	}
	return err
}
