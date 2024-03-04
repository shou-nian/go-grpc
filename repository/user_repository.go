package repository

import (
	"errors"
	"github/riny/go-grpc/user-system/model"
	"golang.org/x/net/context"
	"time"
)

type UserRepositoryManagement interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUserInfo(updateModel interface{}) (interface{}, error)
	QueryUserInfo(ctx context.Context, args any) (*model.User, error)
	QueryUserInfoByToken(ctx context.Context, token string) (*model.User, error)
	CheckEmailIsExisting(email string) (bool, error)
	CheckValidToken(ctx context.Context, token string) (bool, error)
}

type UserRepo struct {
	db *Database
}

var _ UserRepositoryManagement = &UserRepo{}

func NewUserRepository(db *Database) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return ur.db.connection.Create(user).Error
}

func (ur *UserRepo) UpdateUserInfo(updateModel interface{}) (interface{}, error) {
	switch updateModel.(type) {
	case *model.User:
		var user = &model.User{}
		if err := ur.db.connection.Save(updateModel).First(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	case *model.Token:
		var token = &model.Token{}
		token.InDate = time.Now().Add(30 * 24 * time.Hour)
		if err := ur.db.connection.Save(token).First(token).Error; err != nil {
			return nil, err
		}
		return token, nil
	}

	return nil, errors.New("unsupported model")
}

func (ur *UserRepo) QueryUserInfo(ctx context.Context, args any) (*model.User, error) {
	var user = new(model.User)
	switch args.(type) {
	case string:
		if err := ur.db.connection.Model(&model.User{}).Where("email = ?", args).First(user).Error; err != nil {
			return nil, err
		}
	case int32:
		if err := ur.db.connection.Model(&model.User{}).Where("id = ?", args).First(user).Error; err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	return user, nil
}

func (ur *UserRepo) QueryUserInfoByToken(ctx context.Context, token string) (*model.User, error) {
	user := &model.User{}

	if err := ur.db.connection.Where("id = (?)", ur.db.connection.Table("main_token").Select("user_id").Where("token = ?", token)).First(user).Error; err != nil {
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
	t := &model.Token{}
	if err := ur.db.connection.Model(t).Where("token = ?", token).First(t).Error; err != nil {
		return false, err
	}

	if t.InDate.Sub(time.Now()) < 0 {
		return false, nil
	}

	return true, nil
}
