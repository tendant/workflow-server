package dsl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"
	"go.temporal.io/sdk/workflow"
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
	log.Info().Msg("Step 11111")
	err := workflow.ExecuteActivity(ctx, ApprovalActivity, args).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	switch result {
	case "Declined":
		log.Info().Msgf("Workflow completed. ExpenseStatus: %s.\n", result)
		return "", nil
	case "Approved":
		log.Info().Msgf("Continue Workflow. ExpenseStatus: %s.\n", result)
	default:
		// Error
		log.Warn().Msgf("Incorrect status of ApprovelActivity: %s.\n", result)
	}

	// step, wait for the expense report to be approved (or rejected)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	ctx2 := workflow.WithActivityOptions(ctx, ao)
	var status string
	// Step 2
	log.Info().Msg("Step 2222")
	err = workflow.ExecuteActivity(ctx2, ApprovalActivity, args).Get(ctx2, &status)
	if err != nil {
		return "", err
	}
	switch status {
	case "Declined":
		log.Info().Msgf("Workflow completed. ExpenseStatus: %s.\n", status)
		return "", nil
	case "Approved":
		log.Info().Msgf("Continue Workflow. ExpenseStatus: %s.\n", status)
	default:
		// Error
		log.Warn().Msgf("Incorrect status of ApprovelActivity: %s.\n", status)
	}

	return "Completed", nil

}
