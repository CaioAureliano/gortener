package service

import (
	"encoding/json"
	"log"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: r,
	}
}

func (u *UserService) Create(req *model.UserCreateRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var user *model.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		return err
	}

	// TODO: Validate user

	if err := u.userRepository.Create(user); err != nil {
		log.Printf("error to creata a new user: %s", err.Error())
		return err
	}

	return nil
}

func (u *UserService) Exists(req *model.AuthRequest) (bool, error) {
	return u.userRepository.ExistsByEmail(req.Email)
}
