package usecase

import "go.uber.org/zap"

type redirectRepository interface {
	Find(code string) (string, error)
}

type UserRedirect struct {
	repository redirectRepository
	log        *zap.Logger
}

func NewUserRedirect(repository redirectRepository, log *zap.Logger) *UserRedirect {
	return &UserRedirect{repository: repository, log: log}
}

func (u *UserRedirect) Find(code string) (string, error) {
	return u.repository.Find(code)
}
