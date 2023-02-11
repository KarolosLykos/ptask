package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/ptask/domain"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

type ListQueryParams struct {
	Period   *domain.Period
	Timezone *time.Location
	T1       time.Time
	T2       time.Time
}

func GetListQueryParams(
	ctx context.Context,
	logger logger.Logger,
	period, tz, t1, t2 string,
) (*ListQueryParams, error) {
	logger.Trace(ctx, "utils.GetListQueryParams")
	defer logger.Trace(ctx, "utils.GetListQueryParams")

	p, err := parsePeriod(ctx, logger, period)
	if err != nil {
		return nil, err
	}

	timeLoc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", httperrors.ErrInvalidTimezone, err)
	}

	startPoint, err := time.Parse(constants.TimestampLayout, t1)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", httperrors.ErrInvalidStartPoint, err)
	}

	endPoint, err := time.Parse(constants.TimestampLayout, t2)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", httperrors.ErrInvalidEndPoint, err)
	}

	return &ListQueryParams{
		Period:   p,
		Timezone: timeLoc,
		T1:       startPoint.In(timeLoc),
		T2:       endPoint.In(timeLoc),
	}, nil
}

func parsePeriod(ctx context.Context, logger logger.Logger, period string) (*domain.Period, error) {
	logger.Trace(ctx, "utils.parsePeriod")
	defer logger.Trace(ctx, "utils.parsePeriod")

	n := strings.IndexFunc(period, unicode.IsLetter)
	if n == -1 {
		return nil, httperrors.ErrInvalidPeriod
	}

	v, err := strconv.Atoi(period[:n])
	if err != nil {
		return nil, fmt.Errorf("%w:%v", httperrors.ErrInvalidPeriod, err)
	}

	p := strings.ToLower(period[n:])
	switch p {
	case constants.Year, constants.Month, constants.Day, constants.Hour:
		return &domain.Period{Value: v, PeriodType: p}, nil
	default:
		return nil, fmt.Errorf("%w:%v", httperrors.ErrInvalidPeriod, err)
	}
}
