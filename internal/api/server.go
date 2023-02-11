package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/KarolosLykos/ptask/docs"
	"github.com/KarolosLykos/ptask/internal/api/middlewares"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/ptask"
	taskHttp "github.com/KarolosLykos/ptask/internal/ptask/http"
)

type API struct {
	logger  logger.Logger
	addr    string
	handler http.Handler
	server  *http.Server
}

func New(logger logger.Logger, addr string, useCase ptask.UseCase) *API {
	// setting up cors options.
	corsOptions := []handlers.CORSOption{
		handlers.AllowedMethods([]string{http.MethodGet}),
		handlers.AllowedHeaders([]string{"content-type"}),
	}

	router := mux.NewRouter().StrictSlash(true)

	// setting up middlewares.
	m := middlewares.New(logger)

	router.Use(m.RecoverPanic)
	router.Use(m.LogInfo)

	// init task handler.
	h := taskHttp.NewTaskHandler(logger, useCase)

	// setup task routes.
	router = taskHttp.Routes(router, h)

	// setup swagger route.
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", addr)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	)).Methods(http.MethodGet)

	// apply CORS middleware.
	handler := handlers.CORS(corsOptions...)(router)

	return &API{logger: logger, addr: addr, handler: handler}
}

func (a *API) Start(ctx context.Context) {
	a.logger.Info(ctx, "starting server...")

	a.server = &http.Server{
		Addr:         a.addr,
		Handler:      a.handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.logger.Panic(ctx, err, "could not start http server on: ", a.addr)
			}
		}
	}()

	a.logger.Info(ctx, "server started on: ", a.addr)
}

func (a *API) Shutdown(ctx context.Context) {
	a.logger.Debug(ctx, "shutting down server...")

	_ = a.server.Shutdown(ctx)
}
