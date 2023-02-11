package response

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

func TestResponse_Success(t *testing.T) {
	tt := []struct {
		name     string
		status   int
		payload  interface{}
		response string
	}{
		{name: "OK", status: http.StatusOK, payload: []string{
			"20210715T120000Z",
			"20210715T130000Z",
		}, response: "{\"status\":\"success\",\"data\":[\"20210715T120000Z\",\"20210715T130000Z\"]}"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Success(w, tc.status, tc.payload)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			require.Equal(t, tc.status, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.response, string(body))
		})
	}
}

func TestResponse_Error(t *testing.T) {
	tt := []struct {
		name     string
		status   int
		err      error
		payload  interface{}
		response string
	}{
		{name: "wrapped error", status: http.StatusInternalServerError, err: fmt.Errorf("%w:%v", httperrors.ErrInternalServer, errors.New("wrapped error"))},
		{name: "default", status: http.StatusInternalServerError, err: errors.New("new error ")},
		{name: "internal server error", status: http.StatusInternalServerError, err: httperrors.ErrInternalServer},
		{name: "invalid params", status: http.StatusBadRequest, err: httperrors.ErrInvalidTimezone},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Error(w, tc.err)

			resp := w.Result()

			require.Equal(t, tc.status, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		})
	}
}
