package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/http"
	"github.com/cadicallegari/user/mem"
	"github.com/cadicallegari/user/mysql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql/xmysql"
	"github.com/cadicallegari/user/pkg/xhttp"
	"github.com/cadicallegari/user/pkg/xlogger"
	"github.com/cadicallegari/user/pkg/xsignal"
)

var (
	// These var will be injected by the compiler.
	tag       string
	gitCommit string
	buildTime string
)

var cfg struct {
	Logger xlogger.Config     `envconfig:"LOG"`
	HTTP   xhttp.ServerConfig `envconfig:"HTTP"`
	MySQL  xmysql.Config      `envconfig:"MYSQL"`

	PasswordGenerationCost int `envconfig:"PASSWORD_GENERATION_COST" default:"14"`
}

func main() {
	envconfig.MustProcess("user", &cfg)

	log := xlogger.New(&cfg.Logger).WithFields(
		logrus.Fields{"tag": tag, "git_commit": gitCommit},
	)

	cfg.MySQL.Logger = log

	db, err := xmysql.Connect(&cfg.MySQL)
	if err != nil {
		log.WithError(err).
			Error("unable to connect to database")
		return
	}
	defer db.Close()

	storage := mysql.NewStorage(db)
	eventSvc := mem.NewEventService()

	userSrv := user.NewService(storage, eventSvc, cfg.PasswordGenerationCost)

	r := xhttp.NewRouter(log)
	r.Route("/", func(r chi.Router) {
		http.NewUserHandler(r, userSrv)
	})

	httpSrv := xhttp.NewServer(
		&cfg.HTTP,
		xhttp.WithLogger(log.WriterLevel(logrus.ErrorLevel)),
		xhttp.WithRouter(r),
	)
	log.Infof("running http server on: %s", httpSrv.Addr)

	err = httpSrv.ListenAndServe()
	if err != nil {
		log.WithError(err).Error("unable to run http server")
		return
	}
	defer httpSrv.Shutdown()

	sig := xsignal.WaitStopSignal()
	log.WithField("signal", sig).Warn("signal received")
}
