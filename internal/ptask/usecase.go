package ptask

import (
	"context"

	"github.com/KarolosLykos/ptask/internal/ptask/domain"
	"github.com/KarolosLykos/ptask/internal/utils"
)

type UseCase interface {
	GetList(ctx context.Context, params *utils.ListQueryParams) (domain.PtList, error)
}
