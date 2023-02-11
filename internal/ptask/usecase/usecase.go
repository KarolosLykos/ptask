package usecase

import (
	"context"
	"time"

	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger"
	"github.com/KarolosLykos/ptask/internal/ptask"
	"github.com/KarolosLykos/ptask/internal/ptask/domain"
	"github.com/KarolosLykos/ptask/internal/utils"
	"github.com/KarolosLykos/ptask/internal/utils/httperrors"
)

type periodicTaskUC struct {
	logger logger.Logger
}

func NewPeriodicTaskUC(logger logger.Logger) ptask.UseCase {
	return &periodicTaskUC{logger: logger}
}

func (p *periodicTaskUC) GetList(ctx context.Context, params *utils.ListQueryParams) (domain.PtList, error) {
	p.logger.Trace(ctx, "periodicTaskU.GetList")
	defer p.logger.Trace(ctx, "periodicTaskU.GetList")

	task, err := domain.NewPeriodicTask(ctx, p.logger, params.Period, params.Timezone, params.T1)
	if err != nil {
		return nil, err
	}

	list := domain.PtList{}
	for point := task.InvocationPoint; point.Before(params.T2); {
		list = append(list, point.UTC().Format(constants.TimestampLayout))

		switch task.Period.PeriodType {
		case constants.Year:
			point = point.AddDate(task.Period.Value, 0, 0)
		case constants.Month:
			point = point.AddDate(0, task.Period.Value, 0)
		case constants.Day:
			point = point.AddDate(0, 0, task.Period.Value)
		case constants.Hour:
			point = point.Add(time.Duration(task.Period.Value) * time.Hour)
		default:
			return nil, httperrors.ErrInvalidPeriod
		}
	}

	return list, nil
}
