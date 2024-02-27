package repository

import (
	"github/riny/go-grpc/user-system/model"
	"golang.org/x/net/context"
)

type UserRepositoryManagement interface {
	CreateUser(ctx context.Context, user *model.User) error
	QueryUserInfo(ctx context.Context, user *model.User) error
	CheckEmailIsExisting(email string) (bool, error)
	CheckValidToken(ctx context.Context, token string) (bool, error)
}

type UserRepo struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return ur.db.connection.Create(user).Error
}

func (ur *UserRepo) QueryUserInfo(ctx context.Context, email string) (*model.User, error) {
	var user = new(model.User)
	if err := ur.db.connection.Model(&model.User{}).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) CheckEmailIsExisting(email string) (bool, error) {
	var count int64
	if err := ur.db.connection.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return true, err
	}

	return count > 0, nil
}

func (ur *UserRepo) CheckValidToken(ctx context.Context, token string) (bool, error) {
	var count int64
	if err := ur.db.connection.Model(&model.User{}).Where("token = ?", token).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
