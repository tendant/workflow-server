package handler

import (
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux, handle Handle) {
	r.Group(func(r chi.Router) {
		r.Route("/api/v2", func(r chi.Router) {
			r.Get("/", handle.Hello)
			r.Get("/expenses", handle.ListExpenses)
			r.With(httpin.NewInput(ExpenseInput{})).Get("/expenses/{eid}", handle.GetExpense)
			r.With(httpin.NewInput(ExpenseInput{})).Post("/expenses/{eid}/approve", handle.ApproveExpense)
			r.With(httpin.NewInput(ExpenseInput{})).Post("/expenses/{eid}/decline", handle.DeclineExpense)
		})
	})
}
