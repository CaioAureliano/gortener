package handler

import "github.com/CaioAureliano/gortener/internal/auth/model"

type userServiceMock struct {
}

func (u userServiceMock) Create(req *model.UserCreateRequest) error {
	return nil
}

func (u userServiceMock) Exists(email string) (bool, error) {
	return false, nil
}

func (u userServiceMock) GetByField(value, key string) (*model.User, error) {
	return nil, nil
}
