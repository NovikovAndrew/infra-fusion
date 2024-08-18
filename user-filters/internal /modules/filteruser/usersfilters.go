package filteruser

import (
	"context"
	"log/slog"

	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /models"
	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /types"
)

type Store interface {
	GetFilter(ctx context.Context, filterID types.FilterID) (*models.Filter, error)
	GetFiltersByUserID(ctx context.Context, userID types.USerID) ([]*models.Filter, error)
	SaveUserFilter(ctx context.Context, userFilter *models.UserFilter, fn func(ctx context.Context) error) error
}

type DataBus interface {
	ProduceFilter(ctx context.Context, filter *models.Filter) error
}

type UserFilter struct {
	ctx     context.Context
	logger  *slog.Logger
	store   Store
	dataBus DataBus
}

func NewUserFilter(ctx context.Context, logger *slog.Logger, store Store, dataBus DataBus) *UserFilter {
	return &UserFilter{
		ctx:     ctx,
		logger:  logger,
		store:   store,
		dataBus: dataBus,
	}
}

func (uf *UserFilter) GetFilter(ctx context.Context, filterID types.FilterID) (*models.Filter, error) {
	return uf.store.GetFilter(ctx, filterID)
}

func (uf *UserFilter) GetFiltersByUserID(ctx context.Context, userID types.USerID) ([]*models.Filter, error) {
	return uf.store.GetFiltersByUserID(ctx, userID)
}

func (uf *UserFilter) SaveUserFilter(ctx context.Context, userFilter *models.UserFilter) error {
	execFunc := func(ctx context.Context) error {
		return uf.dataBus.ProduceFilter(ctx, convertFromUserFilterToFilter(userFilter))
	}

	return uf.store.SaveUserFilter(ctx, userFilter, execFunc)
}

func convertFromUserFilterToFilter(userFilter *models.UserFilter) *models.Filter {
	return &models.Filter{
		FilterID: userFilter.FilterID,
	}
}
