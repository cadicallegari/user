package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cadicallegari/user"
)

// UserHandler implements http interface
type UserHandler struct {
	userSvc user.Service
}

// NewUserHandler set the http routes for user
func NewUserHandler(r chi.Router, userSvc user.Service) *UserHandler {
	h := &UserHandler{
		userSvc: userSvc,
	}

	r.Get("/v1/users", h.get)
	r.Post("/v1/users", h.create)
	r.Route("/v1/users/{id}", func(r chi.Router) {
		// r = r.With(h.loadUser)
		// r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return h
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {

}
