package handlers

import (
	"net/http"
)

type Service interface {
}

type Handler struct {
	services Service
}

func NewHandler(service Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return mux
}
