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
	ExpenseId string
}

func ParseWorkflow(filePath string) (*model.Workflow, error) {
	workflow, err := parser.FromFile(filePath)
	if err != nil {
		return nil, err
	}
	return workflow, nil
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

func DSLWorkflow(ctx workflow.Context, args DSLWorkflowArgs) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	// Step 1
	slog.Info("Step 11111")
	err := workflow.ExecuteActivity(ctx, ApprovalActivity, args).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	switch result {
	case "Declined":
		slog.Info("Workflow completed.", "ExpenseStatus", result)
		return "", nil
	case "Approved":
		slog.Info("Continue Workflow.", "ExpenseStatus", result)
	default:
		// Error
		slog.Warn("Incorrect status of", "ApprovelActivity", result)
	}

	// step, wait for the expense report to be approved (or rejected)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	ctx2 := workflow.WithActivityOptions(ctx, ao)
	var status string
	// Step 2
	slog.Info("Step 2222")
	err = workflow.ExecuteActivity(ctx2, ApprovalActivity, args).Get(ctx2, &status)
	if err != nil {
		return "", err
	}
	switch status {
	case "Declined":
		slog.Info("Workflow completed.", "ExpenseStatus", status)
		return "", nil
	case "Approved":
		slog.Info("Continue Workflow.", "ExpenseStatus", status)
	default:
		// Error
		slog.Warn("Incorrect status of ApprovelActivity.", "ApprovelActivity", status)
	}

	return "Completed", nil

}
