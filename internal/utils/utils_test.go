package utils

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/logger/log"
	"github.com/KarolosLykos/ptask/internal/ptask/domain"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

func TestGetListQueryParams(t *testing.T) {
	l := getLogger()
	ctx := context.TODO()

	tt := []struct {
		name               string
		period, tz, t1, t2 string
		params             *ListQueryParams
		err                error
	}{
		{
			name: "invalid period", period: "a", err: httperrors.ErrInvalidPeriod,
		},
		{
			name: "invalid location", period: "1h", tz: "WrontTZ", err: httperrors.ErrInvalidTimezone,
		},
		{
			name: "invalid start point", period: "1h", tz: "Europe/Athens", t1: "wrong", err: httperrors.ErrInvalidStartPoint,
		},
		{
			name: "invalid end point", period: "1h", tz: "Europe/Athens", t1: "20060102T150405Z", t2: "wrong", err: httperrors.ErrInvalidEndPoint,
		},
		{
			name:   "ok",
			period: "1h",
			tz:     "",
			t1:     "20060102T150405Z",
			t2:     "20060102T150405Z",
			params: &ListQueryParams{
				Period:   &domain.Period{Value: 1, PeriodType: constants.Hour},
				Timezone: time.UTC,
				T1:       time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				T2:       time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := GetListQueryParams(ctx, l, tc.period, tc.tz, tc.t1, tc.t2)
			if err != nil && tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.params, p)
			}
		})
	}
}

func TestParsePeriod(t *testing.T) {
	l := getLogger()
	ctx := context.TODO()

	tt := []struct {
		name   string
		period string
		parsed *domain.Period
		err    error
	}{
		{
			name: "empty period", period: "", parsed: nil, err: httperrors.ErrInvalidPeriod,
		},
		{
			name: "invalid period", period: "a", parsed: nil, err: httperrors.ErrInvalidPeriod,
		},
		{
			name: "invalid type", period: "10w", parsed: nil, err: httperrors.ErrInvalidPeriod,
		},
		{
			name: "1y", period: "1y", parsed: &domain.Period{Value: 1, PeriodType: constants.Year},
		},
		{
			name: "1Mo", period: "1Mo", parsed: &domain.Period{Value: 1, PeriodType: constants.Month},
		},
		{
			name: "2D", period: "2D", parsed: &domain.Period{Value: 2, PeriodType: constants.Day},
		},
		{
			name: "1h", period: "1h", parsed: &domain.Period{Value: 1, PeriodType: constants.Hour},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := parsePeriod(ctx, l, tc.period)
			if err != nil && tc.err != nil {
				assert.Error(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.parsed, p)
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
