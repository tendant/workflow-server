package dsl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"
	"go.temporal.io/sdk/workflow"
)

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
		StartToCloseTimeout: time.Minute * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string

	err := workflow.ExecuteActivity(ctx, ApprovalActivity, args).Get(ctx, &result)

	return result, err
}
