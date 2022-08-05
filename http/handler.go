package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/pkg/xhttp"
	"github.com/cadicallegari/user/pkg/xlogger"
)

type UserHandler struct {
	userSrv user.Service
}

func NewUserHandler(r chi.Router, userSvc user.Service) *UserHandler {
	h := &UserHandler{
		userSrv: userSvc,
	}

	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Route("/{id}", func(r chi.Router) {
			r = r.With(h.loadUser)
			r.Get("/", h.Get)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})

	return h
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var usrReq *user.User
	err := json.NewDecoder(r.Body).Decode(&usrReq)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to decode request")
		xhttp.ResponseWithStatus(ctx, w, http.StatusBadRequest, nil)
		return
	}

	u, err := h.userSrv.Save(ctx, usrReq)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to save user")
		xhttp.ResponseWithStatus(ctx, w, http.StatusInternalServerError, nil)
		return
	}

	xhttp.ResponseWithStatus(ctx, w, http.StatusCreated, u)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	opts := user.NewListOptions()
	err := xhttp.DecodeQuery(r, opts)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to save parse request")
		xhttp.ResponseWithStatus(ctx, w, http.StatusBadRequest, nil)
		return
	}

	list, err := h.userSrv.List(ctx, opts)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to fetch users")
		xhttp.ResponseWithStatus(ctx, w, http.StatusInternalServerError, nil)
		return
	}

	xhttp.ResponseWithStatus(ctx, w, http.StatusOK, list)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var usrReq *user.User
	err := json.NewDecoder(r.Body).Decode(&usrReq)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to decode request")
		xhttp.ResponseWithStatus(ctx, w, http.StatusBadRequest, nil)
		return
	}

	u, err := h.userSrv.Save(ctx, usrReq)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to update user")
		xhttp.ResponseWithStatus(ctx, w, http.StatusInternalServerError, nil)
		return
	}

	xhttp.ResponseWithStatus(ctx, w, http.StatusOK, u)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usr := ctx.Value(userCtxKey).(*user.User)

	xhttp.ResponseWithStatus(ctx, w, http.StatusOK, usr)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usr := ctx.Value(userCtxKey).(*user.User)

	err := h.userSrv.Delete(ctx, usr)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error("unable to delete user")
		xhttp.ResponseWithStatus(ctx, w, http.StatusInternalServerError, nil)
		return
	}

	xhttp.ResponseWithStatus(ctx, w, http.StatusOK, nil)
}

type contextKey string

var userCtxKey = contextKey("user")

func (h *UserHandler) loadUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := xhttp.URLParam(r, "id")

		u, err := h.userSrv.Get(ctx, id)
		if err != nil {
			xlogger.Logger(ctx).WithError(err).Error("unable to fetch user")
			xhttp.ResponseWithStatus(ctx, w, http.StatusInternalServerError, nil)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
