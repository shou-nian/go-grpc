package service

import (
	"context"
	"github.com/riny/go-grpc/user-system/model"
	"github.com/riny/go-grpc/user-system/repository"
	"github.com/riny/go-grpc/user-system/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"regexp"
	"sync"
)

type ImplementedUserServiceServer struct {
	sync.Mutex
	sync.WaitGroup
	UserServiceServer
	UserRepo *repository.UserRepo
	// refresh token
	// delete
}

var _ UserServiceServer = &ImplementedUserServiceServer{}

func NewImplementedUserServiceServer(userRepo *repository.UserRepo) *ImplementedUserServiceServer {
	return &ImplementedUserServiceServer{UserRepo: userRepo}
}

func (s *ImplementedUserServiceServer) Login(ctx context.Context, req *Login) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password is required.")
	}

	user, err := s.UserRepo.QueryUserInfo(nil, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "email or password error. please try again.")
	}

	if req.Password != user.Password {
		return nil, status.Errorf(codes.InvalidArgument, "email or password error. please try again.")
	}

	// generate new token
	t := util.GenerateStrToken(req.Email, req.Password)
	s.Lock()
	token, err := s.UserRepo.QueryTokenInfo(nil, user.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	token.Token = t

	_, err = s.UserRepo.UpdateModel(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	s.Unlock()
	res := &LoginResponse{Id: user.Id, Token: t}
	return res, nil
}

func (s *ImplementedUserServiceServer) Register(ctx context.Context, req *Register) (*RegisterResponse, error) {
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email or password is required. please try again.")
	}

	// check email is valid
	if ok, err := regexp.MatchString(`(?m)^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`, req.Email); !ok {
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		return nil, status.Errorf(codes.InvalidArgument, "email or password is invalid. please try again.")
	}

	// check password is valid
	if ok, err := util.CheckValidPassword(req.Password); !ok {
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		return nil, status.Errorf(codes.InvalidArgument, "email or password is valid. please try again.")
	}

	// check password is consistent
	if req.ConfirmPassword != req.Password {
		return nil, status.Errorf(codes.InvalidArgument, "password and confirm password not equal. please try again.")
	}

	// check email is existing
	if exists, err := s.UserRepo.CheckEmailIsExisting(req.Email); exists {
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		return nil, status.Errorf(codes.InvalidArgument, "your email resgitered. please check again.")
	}

	// create user
	user := &model.User{Email: req.Email, Password: req.Password, VerifyCode: util.GenerateVerifyCode(req.Password)}
	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Task canceled due to timeout or cancellation")
	default:
		s.Lock()
		err := s.UserRepo.CreateUser(nil, user)
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		// update token
		t := util.GenerateStrToken(user.Email, user.Password)
		token := &model.Token{UserId: user.Id, Token: t}
		// first generate token
		_, err = s.UserRepo.UpdateModel(token)
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		s.Unlock()
		res := &RegisterResponse{
			Id:         user.Id,
			Token:      t,
			VerifyCode: user.VerifyCode,
		}
		return res, nil
	}
}

func (s *ImplementedUserServiceServer) QueryUserInfo(ctx context.Context, req *Query) (*QueryResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	auth := md.Get("Authorization")
	user, err := s.UserRepo.QueryUserInfoByToken(nil, auth[0][7:])
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	if req.Id != user.Id && req.Id != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user id and auth not eque. please try again.")
	}

	res := &QueryResponse{Id: user.Id, Email: user.Email}
	return res, nil
}

func (s *ImplementedUserServiceServer) Update(ctx context.Context, req *Update) (*UpdateResponse, error) {
	updatePassword := req.Password
	if ok, err := util.CheckValidPassword(updatePassword); !ok || updatePassword != req.ConfirmPassword {
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		return nil, status.Errorf(codes.InvalidArgument, "password or confirm password is invalid. please check again.")
	}

	md, _ := metadata.FromIncomingContext(ctx)
	auth := md.Get("Authorization")
	user, err := s.UserRepo.QueryUserInfoByToken(nil, auth[0][7:])
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	// update password
	s.Lock()
	user.Password = updatePassword
	updated1, err := s.UserRepo.UpdateModel(user)
	s.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	updatedUser := updated1.(*model.User)

	// update token
	s.Lock()
	token, err := s.UserRepo.QueryTokenInfo(nil, updatedUser.Id)
	if err == nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	token.Token = util.GenerateStrToken(updatedUser.Email, updatedUser.Password)
	updated2, err := s.UserRepo.UpdateModel(token)
	if err == nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	s.Unlock()

	updatedToken := updated2.(*model.Token)

	var res = &UpdateResponse{
		Id:         updatedToken.UserId,
		Token:      updatedToken.Token,
		VerifyCode: updatedUser.VerifyCode,
	}

	return res, nil
}
