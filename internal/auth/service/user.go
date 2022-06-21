package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (u *UserService) GetByField(value, key string) (*model.User, error) {
	user, err := u.userRepository.GetByField(value, key)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserService) Create(req *model.UserCreateRequest) error {
	if err := isEqualsPassword(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := mapRequestToModel(req)

	if exists, _ := u.Exists(req.Email); exists {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("User already exists with email: %s", req.Email))
	}

	if err := encryptPassword(user); err != nil {
		log.Printf("error to encrypt password: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if err := u.userRepository.Create(user); err != nil {
		log.Printf("error to creata a new user: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (u *UserService) Exists(email string) (bool, error) {
	return u.userRepository.ExistsByEmail(email)
}

func mapRequestToModel(req *model.UserCreateRequest) *model.User {
	return &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func encryptPassword(user *model.User) error {
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(encryptPass)
	return nil
}

func isEqualsPassword(req *model.UserCreateRequest) error {
	if req.Password != req.RepeatPassword {
		return errors.New("passwords should be equals")
	}
	return nil
}
