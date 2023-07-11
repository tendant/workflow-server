package dsl

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/activity"
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
		log.Error().Err(err).Msg("Failed to dial Temporal server")
		return nil, err
	}
	return &WorkflowClient{workflowClient, log}, nil
}

func (w *WorkflowClient) Close() {
	w.wfc.Close()
}

func (w *WorkflowClient) GetWorkflowRunAct(namespace, workflowId string) (WorkflowRunAct, error) {
	we := w.wfc.GetWorkflow(context.Background(), workflowId, "")
	runId := we.GetRunID()
	activityId := activity.GetInfo(context.Background()).ActivityID
	var state string
	err := we.Get(context.Background(), &state)
	if err != nil {
		w.log.Error().Err(err).Msg("unable to get Workflow result")
		return WorkflowRunAct{}, err
	}

	runact := WorkflowRunAct{
		Namespace:  namespace,
		WorkflowId: workflowId,
		RunId:      runId,
		ActivityId: activityId,
		State:      state,
		Err:        nil,
	}
	return runact, nil
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
