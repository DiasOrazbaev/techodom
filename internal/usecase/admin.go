package usecase

import (
	"go.uber.org/zap"
	"techodom/internal/entity"
)

type adminRepo interface {
	FindByID(id string) (*entity.Redirect, error)
	All(page int, perPage int) ([]*entity.Redirect, error)
	Create(old, new string) error
	Update(old, new string, id int) error
	Delete(id int) error
}

type Admin struct {
	logger *zap.Logger
	repo   adminRepo
}

func NewAdmin(usecase adminRepo, logger *zap.Logger) *Admin {
	return &Admin{logger: logger, repo: usecase}
}

func (a *Admin) FindByID(id string) (*entity.Redirect, error) {
	return a.repo.FindByID(id)
}

func (a *Admin) All(page int, perPage int) ([]*entity.Redirect, error) {
	return a.repo.All(page, perPage)
}

func (a *Admin) Create(old, new string) error {
	return a.repo.Create(old, new)
}

func (a *Admin) Update(old, new string, id int) error {
	return a.repo.Update(old, new, id)
}

func (a *Admin) Delete(id int) error {
	return a.repo.Delete(id)
}
