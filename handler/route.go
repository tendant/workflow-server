package handler

import (
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux, handle Handle) {
	r.Group(func(r chi.Router) {
		r.Route("/api/v2", func(r chi.Router) {
			r.Get("/", handle.Hello)
			r.Get("/expenses", handle.Hello)
		})
	})
}
