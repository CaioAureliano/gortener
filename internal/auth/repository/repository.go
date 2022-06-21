package repository

import (
	"github.com/CaioAureliano/gortener/internal/auth/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByField(value, field string) (*model.User, error)
	ExistsByEmail(email string) (bool, error)
}
