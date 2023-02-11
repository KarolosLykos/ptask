package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
	"github.com/KarolosLykos/ptask/internal/utils/response"
)

type middleware struct {
	logger logger.Logger
}

func New(logger logger.Logger) *middleware {
	return &middleware{logger: logger}
}

// RecoverPanic middleware handle any panic that may occur.
func (m *middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer func() {
			if err := recover(); err != nil {
				m.logger.Error(ctx, fmt.Errorf("%w: %v", httperrors.ErrRecoverPanic, err), "middleware recovering from panic error")
				response.Error(w, httperrors.ErrRecoverPanic)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// LogInfo logs request information.
func (m *middleware) LogInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !strings.Contains(r.RequestURI, "/swagger/") {
			start := time.Now()
			defer m.logger.Info(ctx, fmt.Sprintf("%s %s took %s", r.Method, r.RequestURI, time.Since(start).String()))
		}

		next.ServeHTTP(w, r)
	})
}
