package postgres

import (
	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /models"
	"github.com/NovikovAndrew/infra-fusion/user-filters/internal /types"
)

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) GetFilterByID(filterID types.FilterID) (*models.Filter, error) {
	panic("implement me")
}
