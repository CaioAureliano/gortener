package service

import (
	"encoding/json"
	"log"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

var userRepository = repository.NewUserRepository()

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

	encryptPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	if err != nil {
		return err
	}

	user.Password = string(encryptPass)

	if err := userRepository.Create(user); err != nil {
		log.Printf("error to creata a new user: %s", err.Error())
		return err
	}

	return nil
}

func (u *UserService) Exists(req *model.AuthRequest) (bool, error) {
	return userRepository.ExistsByEmail(req.Email)
}
