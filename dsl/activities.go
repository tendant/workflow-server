package dsl

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/activity"
)

func ApprovalActivity(ctx context.Context, args DSLWorkflowArgs) (string, error) {
	activityInfo := activity.GetInfo(ctx)
	log.Debug().Any("task_token", activityInfo.TaskToken).Msg("task_token:")
	// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
	// activity will be completed asynchronously when Client.CompleteActivity() is called.
	return "", activity.ErrResultPending
	// return "", errors.New("Bad Request")
	// return "", err
	// return "", fmt.Errorf("register callback failed status:%s", status)
}
