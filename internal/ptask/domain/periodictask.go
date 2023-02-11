package domain

import (
	"context"
	"time"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

type PeriodicTask struct {
	Period          *Period
	InvocationPoint time.Time
	Timezone        *time.Location
}

type Period struct {
	Value      int
	PeriodType string
}

type PtList []string

func NewPeriodicTask(
	ctx context.Context,
	logger logger.Logger,
	period *Period,
	timezone *time.Location,
	startPoint time.Time,
) (*PeriodicTask, error) {
	logger.Trace(ctx, "periodicTask.NewPeriodicTask")
	defer logger.Trace(ctx, "periodicTask.NewPeriodicTask")

	invocationPoint, err := getInvocationPoint(ctx, logger, period, startPoint)
	if err != nil {
		return nil, err
	}

	return &PeriodicTask{Period: period, InvocationPoint: invocationPoint, Timezone: timezone}, nil
}

func getInvocationPoint(
	ctx context.Context,
	logger logger.Logger,
	period *Period,
	startPoint time.Time,
) (time.Time, error) {
	logger.Trace(ctx, "periodicTask.getInvocationPoint")
	defer logger.Trace(ctx, "periodicTask.getInvocationPoint")

	switch period.PeriodType {
	case constants.Year:
		return time.Date(startPoint.Year()+1, 1, 1, 0, 0, 0, 0, startPoint.Location()), nil
	case constants.Month:
		return time.Date(startPoint.Year(), startPoint.Month()+1, 1, 0, 0, 0, 0, startPoint.Location()), nil
	case constants.Day:
		return time.Date(startPoint.Year(), startPoint.Month(), startPoint.Day()+1, 0, 0, 0, 0, startPoint.Location()), nil
	case constants.Hour:
		return time.Date(
			startPoint.Year(),
			startPoint.Month(),
			startPoint.Day(),
			startPoint.Hour()+1,
			0,
			0,
			0,
			startPoint.Location(),
		), nil
	default:
		return time.Time{}, httperrors.ErrInvalidPeriod
	}
}
