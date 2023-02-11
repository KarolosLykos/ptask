package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := &Response{Status: constants.StatusSuccess}

	if payload != nil {
		res.Data = payload
	}

	p, _ := json.Marshal(res)

	_, _ = w.Write(p)
}

func Error(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	errMsg := err.Error()

	if errU := errors.Unwrap(err); errU != nil {
		errMsg = errU.Error()
	}

	if errors.Is(err, httperrors.ErrInvalidPeriod) ||
		errors.Is(err, httperrors.ErrInvalidTimezone) ||
		errors.Is(err, httperrors.ErrInvalidStartPoint) ||
		errors.Is(err, httperrors.ErrInvalidEndPoint) {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := &Response{Status: constants.StatusError, Error: errMsg}

	p, _ := json.Marshal(res)

	_, _ = w.Write(p)
}
