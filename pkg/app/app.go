package app

import (
	"context"
	"log"
	"time"

	"pastecode/pkg/config"
	"pastecode/pkg/paste"

	"go.uber.org/zap"
)

const (
	readTimeout = 10 * time.Second
)

type Application struct {
	Ctx           context.Context
	CtxCancel     context.CancelFunc
	Sugar         *zap.SugaredLogger
	WebserverConf *config.BackendConfig
	Pastecodes    paste.Pastecodes
}

func (a *Application) NewLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error creating ZAP logger")
	}

	sugar := logger.Sugar()
	a.Sugar = sugar
}

func (a *Application) StopLogger() {
	if err := a.Sugar.Sync(); err != nil {
		a.Sugar.Fatalf("%v", err)
	}
}

func (a *Application) NewContext() {
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	a.Ctx = ctx
	a.CtxCancel = cancel
}

func (a *Application) NewConfig() {
	cfg, err := config.NewBackendConfig()
	if err != nil {
		a.Sugar.Fatalf("Err occured while creating backend config: %v", err)
	}
	a.WebserverConf = cfg
	a.Sugar.Infof("Backend config has been created, host: %s, port: %s", a.WebserverConf.Host, a.WebserverConf.Port)
}

func (a *Application) NewPastecodes() {
	a.Pastecodes = paste.NewPastecodes()
}

func NewApp() *Application {
	app := &Application{}

	app.NewLogger()
	app.NewContext()
	app.NewConfig()
	app.NewPastecodes()

	return app
}
