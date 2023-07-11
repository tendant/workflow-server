package handler

import (
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type Handle struct {
	Slog *slog.Logger
}

type TransactionInput struct {
	TransactionId int `in:"path=txid"`
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
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}
