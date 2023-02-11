package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/logger/log"
	mock_ptask "github.com/KarolosLykos/ptask/internal/ptask/mock"
	"github.com/KarolosLykos/ptask/internal/utils/response"
)

func TestTaskHandler_List(t *testing.T) {
	ctx := context.Background()
	l := getLogger()

	tt := []struct {
		name        string
		useCaseStub func(uc *mock_ptask.MockUseCase)
		method      string
		params      map[string]string
		statusCode  int
		status      string
		list        []string
	}{
		{name: "post not allowed", useCaseStub: func(uc *mock_ptask.MockUseCase) {}, method: http.MethodPost, statusCode: http.StatusMethodNotAllowed},
		{name: "put not allowed", useCaseStub: func(uc *mock_ptask.MockUseCase) {}, method: http.MethodPut, statusCode: http.StatusMethodNotAllowed},
		{
			name:        "invalid params",
			useCaseStub: func(uc *mock_ptask.MockUseCase) {},
			method:      http.MethodGet,
			params:      map[string]string{"period": "wrong", "tz": "Europe/Athens", "t1": "20210728T204603Z", "t2": "20210802T123456Z"},
			statusCode:  http.StatusBadRequest,
			status:      constants.StatusError,
		},
		{
			name: "useCase error",
			useCaseStub: func(uc *mock_ptask.MockUseCase) {
				uc.EXPECT().GetList(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, errors.New("something went wrong"))
			},
			method:     http.MethodGet,
			params:     map[string]string{"period": "1d", "tz": "Europe/Athens", "t1": "20210728T204603Z", "t2": "20210802T123456Z"},
			statusCode: http.StatusInternalServerError,
			status:     constants.StatusError,
		},
		{
			name: "ok",
			useCaseStub: func(uc *mock_ptask.MockUseCase) {
				uc.EXPECT().GetList(gomock.Any(), gomock.Any()).Times(1).
					Return([]string{"20210728T210000Z", "20210729T210000Z", "20210730T210000Z", "20210731T210000Z", "20210801T210000Z"}, nil)
			},
			method:     http.MethodGet,
			params:     map[string]string{"period": "1d", "tz": "Europe/Athens", "t1": "20210728T204603Z", "t2": "20210802T123456Z"},
			statusCode: http.StatusOK,
			status:     constants.StatusSuccess,
			list:       []string{"20210728T210000Z", "20210729T210000Z", "20210730T210000Z", "20210731T210000Z", "20210801T210000Z"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			useCase := mock_ptask.NewMockUseCase(ctrl)

			tc.useCaseStub(useCase)

			h := NewTaskHandler(l, useCase)

			router := mux.NewRouter()
			router.HandleFunc("/ptlist", h.List()).Methods(http.MethodGet)

			srv := httptest.NewServer(router)
			defer srv.Close()

			req, err := http.NewRequestWithContext(ctx, tc.method, fmt.Sprintf("%s/ptlist", srv.URL), nil)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			q := req.URL.Query()
			for k, v := range tc.params {
				q.Add(k, v)
			}

			req.URL.RawQuery = q.Encode()

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			defer res.Body.Close()

			assert.Equal(t, tc.statusCode, res.StatusCode)

			if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusBadRequest {
				resp := &response.Response{}
				err = json.NewDecoder(res.Body).Decode(resp)
				require.NoError(t, err)
				assert.Equal(t, tc.status, resp.Status)

				if tc.list != nil {
					list, ok := resp.Data.([]interface{})
					require.True(t, ok)

					assert.ElementsMatch(t, list, tc.list)
				}
			}
		})
	}
}

func getLogger() logger.Logger {
	l := &logrus.Logger{
		Out:          io.Discard,
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
		Level:        logrus.DebugLevel,
		Formatter:    &logrus.JSONFormatter{},
	}

	return log.New(l)
}
