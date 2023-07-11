package handler

import (
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux, handle Handle) {
	r.Group(func(r chi.Router) {
		r.Route("/api/v2/workflow", func(r chi.Router) {
			r.Get("/", handle.Hello)
			r.Get("/transactions", handle.ListTransactions)
			// FIXME: Start workflow POST /transactions/{txid}, action Start
			// get transaction current workflow state approvers
			r.With(httpin.NewInput(TransactionInput{})).Get("/transactions/{txid}", handle.GetTransaction)
			// approve transaction for given activity id, action Approve or Decline
			r.With(httpin.NewInput(TransactionInput{})).Post("/transactions/{txid}", handle.TransactionApprovalAction)
		})
	})
}
