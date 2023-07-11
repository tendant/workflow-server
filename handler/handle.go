package handler

import (
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type Handle struct {
	Slog *slog.Logger
}

type ExpenseInput struct {
	ExpenseId int `in:"path=eid"`
}

func (h Handle) Hello(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) ListExpenses(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) GetExpense(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) ApproveExpense(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func (h Handle) DeclineExpense(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}
