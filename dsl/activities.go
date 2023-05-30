package dsl

import "context"

func ApprovalActivity(ctx context.Context, args DSLWorkflowArgs) (string, error) {
	return "Approved", nil
}
