package handler

import (
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type Handle struct {
	Slog *slog.Logger
}

func (h Handle) Hello(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}
