package xhttp

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"

	"github.com/cadicallegari/user/pkg/xlogger"
)

var schemaDecoder = schema.NewDecoder()
var DecodeQueryIgnoringUnknownKeys = true

func init() {
	schemaDecoder.IgnoreUnknownKeys(DecodeQueryIgnoringUnknownKeys)
}

type ServerConfig struct {
	Addr string `envconfig:"ADDR" default:"0.0.0.0:80"`
}

func (cfg *ServerConfig) setDefault() {
	if cfg.Addr == "" {
		cfg.Addr = ":http"
	}
}

type Server struct {
	http.Server
}

type serverOption func(*Server)

func NewServer(cfg *ServerConfig, opts ...serverOption) *Server {
	if cfg == nil {
		cfg = new(ServerConfig)
	}
	cfg.setDefault()

	srv := &Server{
		Server: http.Server{
			Addr: cfg.Addr,
		},
	}

	for _, optFn := range opts {
		optFn(srv)
	}

	return srv
}

func WithLogger(w io.Writer) func(*Server) {
	return func(srv *Server) {
		srv.ErrorLog = log.New(w, "", 0)
	}
}

func WithRouter(r *Router) func(*Server) {
	return func(srv *Server) {
		srv.Handler = r
	}
}

func (srv *Server) ListenAndServe() error {
	ch := make(chan error)
	go func() {
		defer close(ch)

		err := srv.Server.ListenAndServe()
		ch <- err
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second):
	}

	return nil
}

var ShutdownTimeout = 5 * time.Second

var ErrServerClosed = http.ErrServerClosed

func (srv *Server) Shutdown() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, ShutdownTimeout)
	defer cancel()

	err := srv.Server.Shutdown(ctx)
	if err != nil && srv.ErrorLog != nil {
		srv.ErrorLog.Printf("unable to shutdown http server: %v", err)
	}
}

type Router struct {
	chi.Router
}

func NewRouter(log *logrus.Entry) *Router {
	r := &Router{
		Router: chi.NewMux(),
	}

	r.Use(middleware.CleanPath)
	r.Use(middleware.Heartbeat("/ping"))

	if log != nil {
		r.Use(loggerMiddleware(log))
	}

	return r
}

func loggerMiddleware(logger *logrus.Entry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := setLoggerRequestFields(logger, r)
			ctx = xlogger.SetLogger(ctx, log)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				status := ww.Status()
				log = setLoggerResponseFields(log, status, ww.BytesWritten(), time.Since(t1))

				if status >= 500 {
					log.Error("http request complete")
					return
				}

				log.Info("http request complete")
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func setLoggerRequestFields(log *logrus.Entry, r *http.Request) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"http_request.method":     r.Method,
		"http_request.user_agent": r.UserAgent(),
		"http_request.path":       r.URL.Path,
	})
}

func setLoggerResponseFields(log *logrus.Entry, status, bytes int, elapsed time.Duration) *logrus.Entry {
	fields := logrus.Fields{
		"http_response.status":       status,
		"http_response.bytes_length": bytes,
	}
	if elapsed > 0 {
		fields["http_response.elapsed"] = elapsed.String()
		fields["http_response.elapsed_ms"] = float64(elapsed.Nanoseconds()) / float64(time.Millisecond)
	}

	return log.WithFields(fields)
}

func DecodeQuery(r *http.Request, v interface{}) error {
	err := schemaDecoder.Decode(v, r.URL.Query())
	if err != nil {
		return err
	}

	return nil
}

func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

var (
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")
	jsonContentType   = "application/json; charset=UTF-8"
)

func ResponseWithStatus(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) {
	if body != nil {
		if ct := w.Header().Get(contentTypeHeader); ct == "" {
			w.Header().Set(contentTypeHeader, jsonContentType)
		}
	}

	w.WriteHeader(statusCode)
	if body == nil {
		return
	}

	if r, ok := body.(io.Reader); ok {
		io.Copy(w, r)
		return
	}

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		xlogger.Logger(ctx).WithError(err).Error(http.StatusText(statusCode))
	}
}
