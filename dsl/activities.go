package dsl

import (
	"context"

	"go.temporal.io/sdk/activity"
	"golang.org/x/exp/slog"
)

func ApprovalActivity(ctx context.Context, args DSLWorkflowArgs) (string, error) {
	activityInfo := activity.GetInfo(ctx)
	slog.Debug("task_token:", "task_token", activityInfo.TaskToken)
	// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
	// activity will be completed asynchronously when Client.CompleteActivity() is called.
	return "", activity.ErrResultPending
	// return "", errors.New("Bad Request")
	// return "", err
	// return "", fmt.Errorf("register callback failed status:%s", status)
}

type TransactionApprovalParams struct {
	Approver string
}

func TransactionApprovalActivity(ctx context.Context, params TransactionApprovalParams) (string, error) {
	activityInfo := activity.GetInfo(ctx)
	slog.Info("activity info:", "activity RunID", activityInfo.WorkflowExecution.RunID, "activityId", activityInfo.ActivityID)

	// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
	// activity will be completed asynchronously when Client.CompleteActivity() is called.
	return "", activity.ErrResultPending
	// return "", errors.New("Bad Request")
	// return "", err
	// return "", fmt.Errorf("register callback failed status:%s", status)
}
