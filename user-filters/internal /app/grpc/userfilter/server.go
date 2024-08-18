package userfilter

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /models"
	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /types"
	"github.com/NovikovAndrew/infra-fusion/user-filters/pkg/pb"
)

type UserFilterModule interface {
	SaveUserFilter(ctx context.Context, userFilter *models.UserFilter)
	GetFilter(ctx context.Context, filterID types.FilterID) (*models.Filter, error)
	GetFiltersByUserID(ctx context.Context, userID types.USerID) ([]*models.Filter, error)
}

type Server struct {
	pb.UnimplementedUserFiltersServer
	logger     *slog.Logger
	userFilter UserFilterModule
}

func (s *Server) Register(register grpc.ServiceRegistrar) {
	pb.RegisterUserFiltersServer(register, s)
}

func NewServer(logger *slog.Logger, userFilter UserFilterModule) *Server {
	return &Server{
		logger:     logger,
		userFilter: userFilter,
	}
}
