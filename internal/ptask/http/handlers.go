package http

import (
	"net/http"

	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/ptask"
	"github.com/KarolosLykos/ptask/internal/utils"
	"github.com/KarolosLykos/ptask/internal/utils/response"
)

type TaskHandler struct {
	logger  logger.Logger
	useCase ptask.UseCase
}

func NewTaskHandler(logger logger.Logger, useCase ptask.UseCase) *TaskHandler {
	return &TaskHandler{
		logger:  logger,
		useCase: useCase,
	}
}

// List returns all matching timestamps of a periodic task
//
//	@Summary		Returns all matching timestamps of a periodic task between 2 points in time.
//	@Accept			json
//	@Produce		json
//	@Param			period	query	string	false	"Period"		example(1y,1mo,1d,1h)
//	@Param			tz		query	string	false	"Timezone"		example(America/Los_Angeles)
//	@Param			t1		query	string	false	"Start point"	example(20060102T150405Z)
//	@Param			t2		query	string	false	"End point"		example(20060102T150405Z)
//	@Success		200
//	@Failure		400
//	@Failure		500
//
//	@Router			/ptlist [get]
func (t *TaskHandler) List() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		period := r.URL.Query().Get("period")
		tz := r.URL.Query().Get("tz")
		t1 := r.URL.Query().Get("t1")
		t2 := r.URL.Query().Get("t2")

		params, err := utils.GetListQueryParams(ctx, t.logger, period, tz, t1, t2)
		if err != nil {
			t.logger.Error(ctx, err, "could not parse query params")
			response.Error(w, err)

			return
		}

		list, err := t.useCase.GetList(ctx, params)
		if err != nil {
			t.logger.Error(ctx, err, "could not get matching task list")
			response.Error(w, err)
		}

		response.Success(w, http.StatusOK, list)
	}
}
