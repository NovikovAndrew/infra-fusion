package userfilter

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /models"
	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /types"
	"github.com/NovikovAndrew/infra-fusion/user-filters/pkg/pb"
)

func (s *Server) GetFiltersByUserID(ctx context.Context, req *pb.GetFiltersByUserIDRequest) (*pb.GetFiltersByUserIDResponse, error) {
	filters, err := s.userFilter.GetFiltersByUserID(ctx, convertGetFiltersByUserIDRequestToUserID(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "[GetFiltersByUserID] failed to get filters by user id: %s", err.Error())
	}

	return &pb.GetFiltersByUserIDResponse{
		UserID:  req.UserID,
		Filters: convertFiltersToPbFilters(filters),
	}, nil
}

func pbConvertGetFiltersByUserIDRequestToUserID(req *pb.GetFiltersByUserIDRequest) types.USerID {
	return types.USerID(req.UserID)
}

func convertFiltersToPbFilters(filters []*models.Filter) []*pb.Filter {
	res := make([]*pb.Filter, len(filters))
	for i := range filters {
		res[i] = &pb.Filter{
			FilterID: filters[i].FilterID,
		}
	}

	return res
}
