package dsl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/exp/slog"
)

const DSLWorkflowTaskQueue = "DSL_WORKFLOW_TASK_QUEUE"

type DSLWorkflowArgs struct {
	Id         string
	Type       string
	EntityType string
	EntityId   string
	DSLStr     string
}

func ParseWorkflow(filePath string) (*model.Workflow, error) {
	workflow, err := parser.FromFile(filePath)
	if err != nil {
		return nil, err
	}
	return workflow, nil
}

func ParseWorkflowDSL(dslStr string) (*model.Workflow, error) {
	workflow, err := parser.FromYAMLSource([]byte(dslStr))
	return workflow, err
}

func GetWorkflowStateByName(name string, workflow *model.Workflow) (model.State, error) {
	// start := workflow.Start.StateName
	if strings.TrimSpace(name) == "" {
		return model.State{}, errors.New("Workflow State name can't be empty")
	}
	for _, s := range workflow.States {
		if s.Name == name {
			return s, nil
		}
	}
	return model.State{}, errors.New(fmt.Sprintf("No State Found for State name: %s", name))
}

func GetStartingWorkflowState(workflow *model.Workflow) (model.State, error) {
	start := workflow.Start.StateName
	return GetWorkflowStateByName(start, workflow)
}

func ExecuteDSLState(ctx workflow.Context, args DSLWorkflowArgs, state model.State) (string, error) {
	slog.Info("State Type", "stateType", state.Type)
	slog.Info("State ActionMode", "actionMode", state.ActionMode)
	slog.Info("State Actions", "actions", state.OperationState.Actions)
	for i, v := range state.OperationState.Actions {
		slog.Info("Runing Action", "i", i, "action", v.Name, "functionRef", v.FunctionRef.RefName)
	}
	return "Completed", nil

}

func ExecuteDSLWorkflow(ctx workflow.Context, args DSLWorkflowArgs, dslWorkflow *model.Workflow) (string, error) {
	slog.Info("Start executing with state name", "stateName", dslWorkflow.Start.StateName)
	startStateName := dslWorkflow.Start.StateName
	state, err := GetWorkflowStateByName(startStateName, dslWorkflow)
	if err != nil {
		slog.Error("Failed getting workflow state by name", "startStateName", startStateName)
		return "", err
	}
	return ExecuteDSLState(ctx, args, state)
}

func DSLWorkflow(ctx workflow.Context, args DSLWorkflowArgs) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	slog.Info("Parsing Workflow DSL")
	dslWorkflow, err := ParseWorkflowDSL(args.DSLStr)
	if err != nil {
		slog.Error("Failed Parse Workflow DSL", "DSLStr", args.DSLStr)
		return "", err
	}

	slog.Info("Start Executing DSL Workflow")
	return ExecuteDSLWorkflow(ctx, args, dslWorkflow)

}
