package service

import (
	"context"
	"errors"
	"github/riny/go-grpc/user-system/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"regexp"
)

type ImplementedUserServiceServer struct {
	UserServiceServer
}

func (s *ImplementedUserServiceServer) Login(ctx context.Context, req *Login) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required.")
	}

	if req.Email == "test@example.com" && req.Password == "admin" {
		resp := &LoginResponse{
			Id:    0,
			Token: "admin",
		}
		return resp, nil
	}

	return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
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

	// create user
	resp, err := createUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return resp, nil
}

func (s *ImplementedUserServiceServer) QueryUserInfo(ctx context.Context, req *Query) (*QueryResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("user id is required")
	}

	return &QueryResponse{Id: req.Id, Email: "temp@temp.com"}, nil
}

func createUser(ctx context.Context, req *Register) (*RegisterResponse, error) {
	resp := &RegisterResponse{}

	return resp, nil
}
