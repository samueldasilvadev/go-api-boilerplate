package user

import (
	domain "go-skeleton/application/domain/dummy"

	"go-skeleton/pkg/database"
)

type Repository struct {
	*database.MySql
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create() error {
	return nil
}

func (r *Repository) FindById(id string) (*domain.Dummy, error) {
	return nil, nil
}
