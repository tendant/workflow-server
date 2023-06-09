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

const (
	DSLWorkflowTaskQueue = "DSL_WORKFLOW_TASK_QUEUE"
	DSLWorkflowQueryType = "state"
)

type DSLWorkflowArgs struct {
	Id         string
	Type       string
	EntityType string
	EntityId   string
	DSLStr     string
}

type DSLWorkflowStateResult struct {
	CurrentState string
	ActionName   string
	Approver     string
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

func ExecuteDSLAction(ctx workflow.Context, args DSLWorkflowArgs, action model.Action, dslState *DSLWorkflowStateResult) (string, error) {
	slog.Info("Runing Action", "name", action.Name, "functionRef", action.FunctionRef.RefName)
	var activityResult string
	params := TransactionApprovalParams{
		Approver: "admin",
	}
	err := workflow.ExecuteActivity(ctx, action.FunctionRef.RefName, params).Get(ctx, &activityResult)
	if err != nil {
		return "", err
	}
	dslState.ActionName = action.Name
	dslState.Approver = params.Approver
	dslState.CurrentState = activityResult
	return activityResult, nil
}

func ExecuteDSLState(ctx workflow.Context, args DSLWorkflowArgs, dslWorkflow *model.Workflow, state model.State, dslState *DSLWorkflowStateResult) (string, error) {
	slog.Info("State Type", "stateType", state.Type)
	slog.Info("State ActionMode", "actionMode", state.ActionMode)
	slog.Info("State Actions", "actions", state.OperationState.Actions)
	stateResult := ""
	for i, v := range state.OperationState.Actions {
		slog.Info("Executing Action", "i", i)
		actionResult, err := ExecuteDSLAction(ctx, args, v, dslState)
		if err != nil {
			slog.Error("Failed Executing Action", "err", err)
			return "", err
		} else {
			if v.ActionDataFilter.UseResults {
				stateResult = actionResult
			}
		}
	}
	if state.End != nil || stateResult == "Declined" {
		return stateResult, nil
	} else {
		nextStateName := state.Transition.NextState
		nextState, err := GetWorkflowStateByName(nextStateName, dslWorkflow)
		if err != nil {
			slog.Error("Failed get workflow state by name", "nextState", nextStateName)
			return "", err
		}
		return ExecuteDSLState(ctx, args, dslWorkflow, nextState, dslState)
	}
}

func ExecuteDSLWorkflow(ctx workflow.Context, args DSLWorkflowArgs, dslWorkflow *model.Workflow, dslState *DSLWorkflowStateResult) (string, error) {
	slog.Info("Start executing with state name", "stateName", dslWorkflow.Start.StateName, "dslState", dslState)
	startStateName := dslWorkflow.Start.StateName
	state, err := GetWorkflowStateByName(startStateName, dslWorkflow)
	if err != nil {
		slog.Error("Failed getting workflow state by name", "startStateName", startStateName)
		return "", err
	}
	return ExecuteDSLState(ctx, args, dslWorkflow, state, dslState)
}

func DSLWorkflow(ctx workflow.Context, args DSLWorkflowArgs) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	slog.Info("Setting up Workflow Query Handler", "QueryType", DSLWorkflowQueryType)
	// setup query handler for query type "state"
	dslState := DSLWorkflowStateResult{
		CurrentState: "Started",
	} // This could be any serializable struct.
	err := workflow.SetQueryHandler(ctx, DSLWorkflowQueryType, func() (DSLWorkflowStateResult, error) {
		return dslState, nil
	})
	if err != nil {
		slog.Error("SetQueryHandler failed: ", "err", err)
		return "", err
	}

	slog.Info("Parsing Workflow DSL")
	dslWorkflow, err := ParseWorkflowDSL(args.DSLStr)
	if err != nil {
		slog.Error("Failed Parse Workflow DSL", "DSLStr", args.DSLStr, "err", err)
		return "", err
	}

	slog.Info("Start Executing DSL Workflow")
	return ExecuteDSLWorkflow(ctx, args, dslWorkflow, &dslState)
}
