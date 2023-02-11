package domain

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/logger/log"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

func TestGetInvocationPoint(t *testing.T) {
	ctx := context.Background()
	l := getLogger()

	tt := []struct {
		name            string
		period          *Period
		start           string
		invocationPoint string
		err             error
	}{
		{
			name: "invalid period",
			period: &Period{
				Value:      1,
				PeriodType: "wrong",
			},
			start: "2014-07-14 23:46:03 +0300 EEST",
			err:   httperrors.ErrInvalidPeriod,
		},
		{
			name: "year period",
			period: &Period{
				Value:      1,
				PeriodType: constants.Year,
			},
			start:           "2014-07-14 23:46:03 +0300 EEST",
			invocationPoint: "2015-01-01 00:00:00 +0200 EET",
		},
		{
			name: "month period",
			period: &Period{
				Value:      1,
				PeriodType: constants.Month,
			},
			start:           "2021-07-14 23:46:03 +0300 EEST",
			invocationPoint: "2021-08-01 00:00:00 +0300 EEST",
		},
		{
			name: "day period",
			period: &Period{
				Value:      1,
				PeriodType: constants.Day,
			},
			start:           "2021-07-28 23:46:03 +0300 EEST",
			invocationPoint: "2021-07-29 00:00:00 +0300 EEST",
		},
		{
			name: "hour period",
			period: &Period{
				Value:      1,
				PeriodType: constants.Hour,
			},
			start:           "2021-07-14 23:46:03 +0300 EEST",
			invocationPoint: "2021-07-15 00:00:00 +0300 EEST",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			format := "2006-01-02 15:04:05 -0700 MST"
			start, err := time.Parse(format, tc.start)
			require.NoError(t, err)

			invocationPoint, err := getInvocationPoint(ctx, l, tc.period, start)
			if err != nil && tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				expected, err := time.Parse(format, tc.invocationPoint)
				require.NoError(t, err)

				assert.Equal(t, expected, invocationPoint)
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
