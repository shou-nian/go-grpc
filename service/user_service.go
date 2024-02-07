package service

import (
	"context"
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

	if req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email or password is required. please try again.")
	}

	if ok, err := regexp.MatchString("^[\\w.-]+@[\\w.-]+\\.\\w+$", req.Email); !ok {
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		return nil, status.Errorf(codes.InvalidArgument, "email or password is valid. please try again.")
	}
	if ok, err := regexp.MatchString("^(?=.*[A-Z])(?=.*[+-=./?])[A-Za-z0-9+-=./?]{8,}$", req.Password); !ok {
		if err != nil {
			return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		return nil, status.Errorf(codes.InvalidArgument, "email or password is valid. please try again.")
	}
	if req.ConfirmPassword != req.Password {
		return nil, status.Errorf(codes.InvalidArgument, "password and confirm password not equal. please try again.")
	}
	// 邮箱是否冲突校验

	// 全部校验通过
	resp := &RegisterResponse{
		Id:         1,
		Token:      "admin",
		VerifyCode: "1",
	}
	return resp, nil
}
