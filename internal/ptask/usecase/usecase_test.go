package usecase

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
	"github.com/KarolosLykos/ptask/internal/ptask/domain"
	"github.com/KarolosLykos/ptask/internal/utils"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

func TestPeriodicTaskUC_GetList(t *testing.T) {
	l := getLogger()
	ctx := context.TODO()

	useCase := NewPeriodicTaskUC(l)

	tt := []struct {
		name                 string
		periodicType, t1, t2 string
		list                 []string
		err                  error
	}{
		{
			name:         "invalid period",
			periodicType: "wrong",
			t1:           "20140714T204603Z",
			t2:           "20210915T123456Z",
			err:          httperrors.ErrInvalidPeriod,
		},
		{
			name:         "1y",
			periodicType: constants.Year,
			t1:           "20140714T204603Z",
			t2:           "20210915T123456Z",
			list:         []string{"20141231T220000Z", "20151231T220000Z", "20161231T220000Z", "20171231T220000Z", "20181231T220000Z", "20191231T220000Z", "20201231T220000Z"},
		},
		{
			name:         "1mo",
			periodicType: constants.Month,
			t1:           "20210714T204603Z",
			t2:           "20210915T123456Z",
			list:         []string{"20210731T210000Z", "20210831T210000Z"},
		},
		{
			name:         "1d",
			periodicType: constants.Day,
			t1:           "20210728T204603Z",
			t2:           "20210802T123456Z",
			list:         []string{"20210728T210000Z", "20210729T210000Z", "20210730T210000Z", "20210731T210000Z", "20210801T210000Z"},
		},
		{
			name:         "1h",
			periodicType: constants.Hour,
			t1:           "20210714T204603Z",
			t2:           "20210715T123456Z",
			list:         []string{"20210714T210000Z", "20210714T220000Z", "20210714T230000Z", "20210715T000000Z", "20210715T010000Z", "20210715T020000Z", "20210715T030000Z", "20210715T040000Z", "20210715T050000Z", "20210715T060000Z", "20210715T070000Z", "20210715T080000Z", "20210715T090000Z", "20210715T100000Z", "20210715T110000Z", "20210715T120000Z"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			params := getParams(t, tc.periodicType, tc.t1, tc.t2)
			list, err := useCase.GetList(ctx, params)
			if err != nil && tc.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, httperrors.ErrInvalidPeriod.Error())
			} else {
				require.NoError(t, err)
				assert.ElementsMatch(t, list, tc.list)
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

func getParams(t *testing.T, periodType, t1, t2 string) *utils.ListQueryParams {
	timezone, err := time.LoadLocation("Europe/Athens")
	require.NoError(t, err)

	start, err := time.Parse(constants.TimestampLayout, t1)
	require.NoError(t, err)
	end, err := time.Parse(constants.TimestampLayout, t2)
	require.NoError(t, err)

	return &utils.ListQueryParams{
		Period:   &domain.Period{Value: 1, PeriodType: periodType},
		Timezone: timezone,
		T1:       start.In(timezone),
		T2:       end.In(timezone),
	}
}
