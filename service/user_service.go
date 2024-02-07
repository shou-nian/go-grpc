package service

import (
	"context"
	"google.golang.org/protobuf/types/known/structpb"
	"net/http"
)

type ImplementedUserServiceServer struct {
	UserServiceServer
}

func (s *ImplementedUserServiceServer) Login(ctx context.Context, req *Login) (*LoginResponse, error) {
	defer ctx.Done()
	resp := &LoginResponse{}

	if req.Email == "" || req.Password == "" {
		str, err := structpb.NewStruct(map[string]interface{}{"msg": "email and password is required"})
		if err != nil {
			return nil, err
		}
		resp.Msg = str
		return resp, nil
	}

	if req.Email == "test@example.com" && req.Password == "admin" {
		resp.Id = 0
		resp.Token = "admin"
		return resp, nil
	}

	msg, err := structpb.NewStruct(map[string]interface{}{"msg": http.StatusText(http.StatusInternalServerError)})
	if err != nil {
		return nil, err
	}
	resp.Msg = msg
	return resp, nil
}
