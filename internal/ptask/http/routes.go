package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/KarolosLykos/ptask/internal/ptask"
)

func Routes(router *mux.Router, taskHandler ptask.Handlers) *mux.Router {
	router.HandleFunc("/ptlist", taskHandler.List()).Methods(http.MethodGet)

	return router
}
