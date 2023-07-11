package handler

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ggicci/httpin"
	"github.com/go-chi/render"
	"github.com/tendant/workflow-server/dsl"
	"go.temporal.io/sdk/client"
	"golang.org/x/exp/slog"
)

type Handle struct {
	Slog   *slog.Logger
	Client client.Client
	Ef     embed.FS
}

type TransactionGetInput struct {
	TransactionId int `in:"path=txnid"`
}

type TransactionPostBody struct {
	TransactionId int    `json:"txnID"`
	Action        string `json:"action"`
	Filename      string `json:"filename"`
}

type TransactionPostInput struct {
	Payload *TransactionPostBody `in:"body=json"`
}

type TransactionPostResponseStart struct {
	WorkflowID    string `json:"workflow_id"`
	WorkflowRunID string `json:"workflow_run_id"`
}

func (h Handle) Hello(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) ListTransactions(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) GetTransaction(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) TransactionApprovalAction(w http.ResponseWriter, r *http.Request) {
	body := r.Context().Value(httpin.Input).(*TransactionPostInput)
	h.Slog.Info("Handling TransactionApprovalAction", "payload", body.Payload)

	txnID := body.Payload.TransactionId
	action := body.Payload.Action
	if txnID < 1 || action == "" {
		http.Error(w, "invalid txnID or action", http.StatusBadRequest)
		return
	}

	filename := body.Payload.Filename
	dslStr, err := h.Ef.ReadFile(fmt.Sprintf("static/%s", filename))
	if err != nil {
		h.Slog.Error("Failed reading workflow DSL from file", "file", filename)
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}
	wfType := "tx"
	entityType := "approval"
	entityId := strconv.FormatInt(int64(txnID), 10)
	id := fmt.Sprintf("%s-%s-%s", wfType, entityType, entityId)

	switch strings.ToLower(action) {
	case "start":
		args := dsl.DSLWorkflowArgs{
			Id:         id,
			Type:       wfType,
			EntityType: entityType,
			EntityId:   entityId,
			DSLStr:     string(dslStr),
		}
		h.Slog.Info("Exeucte workflow", "args", args)
		// Start the Workflow
		options := client.StartWorkflowOptions{
			ID:        id,
			TaskQueue: dsl.DSLWorkflowTaskQueue,
		}
		h.Slog.Info("Workflow options:", "options", options)

		we, err := h.Client.ExecuteWorkflow(context.Background(), options, dsl.DSLWorkflow, args)
		if err != nil {
			h.Slog.Error("Failed to execute Workflow", "err", err)
			http.Error(w, "unable to execute Workflow", http.StatusInternalServerError)
			return
		}
		weID := we.GetID()
		weRunID := we.GetRunID()
		h.Slog.Info("Started Workflow", "workflowID", weID, "runID", weRunID)

		render.JSON(w, r, TransactionPostResponseStart{
			WorkflowID:    weID,
			WorkflowRunID: weRunID,
		})
	default:
		h.Slog.Warn("WIP", "action", action)
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	}
}
