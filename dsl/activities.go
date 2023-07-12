package dsl

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
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
	Approver string `json:"approver,omitempty"`
}

type TransactionPostBody struct {
	TransactionId int            `json:"txnID"`
	Action        string         `json:"action"`
	Approver      string         `json:"approver,omitempty"`
	Activity      WorkflowRunAct `json:"activity,omitempty"`
}

func TransactionApprovalActivity(ctx context.Context, params TransactionApprovalParams) (string, error) {
	activityInfo := activity.GetInfo(ctx)
	// WIP: Send RunID and activityID to external system, so that external system can call Temporal API to complete the Activity
	slog.Info("activity info:", "activity RunID", activityInfo.WorkflowExecution.RunID, "activityId", activityInfo.ActivityID)
	runact := WorkflowRunAct{
		Namespace:  activityInfo.WorkflowNamespace,
		WorkflowId: activityInfo.WorkflowExecution.ID,
		RunId:      activityInfo.WorkflowExecution.RunID,
		ActivityId: activityInfo.ActivityID,
	}
	strs := strings.Split(runact.WorkflowId, "-") // FIXME: use util
	if len(strs) < 3 {
		return "", errors.New("malformed workflowID")
	}
	txnID, err := strconv.Atoi(strs[2])
	if err != nil {
		return "", errors.New("malformed workflowID")
	}
	body := TransactionPostBody{
		TransactionId: txnID,
		Action:        "register",
		Approver:      params.Approver,
		Activity:      runact,
	}
	client := resty.New()
	slog.Info("register activity info...", "body", body)
	resp, err := client.R().
		SetBody(body).
		SetResult(&TransactionPostBody{}). // FIXME: what response struct
		Post("http://localhost:4000/api/v2/workflow/transactions")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("failed to register activity in external system")
	}
	slog.Info("registered activity info", "resp", resp)

	// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
	// activity will be completed asynchronously when Client.CompleteActivity() is called.
	return "", activity.ErrResultPending
}
