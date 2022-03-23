package main

import (
	userRepository "app/internal/user"
	user "app/pb"
	"app/pkg/client/redis"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type UserServiceServer struct {
	repository userRepository.Repository
	cache      redis.Client
}

func (s *UserServiceServer) Create(ctx context.Context, req *user.CreateReq) (*user.CreateRes, error) {
	uid, err := s.repository.Create(context.TODO(), &userRepository.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return &user.CreateRes{
			User: nil,
			Status: &user.Status{
				Code: user.Status_FAILED,
				Msg:  "failed",
			},
		}, err
	}

	return &user.CreateRes{
		User: &user.User{
			Id:       uid,
			Email:    req.Email,
			Password: req.Password,
		},
		Status: &user.Status{
			Code: user.Status_SUCCESS,
			Msg:  "success",
		},
	}, nil
}

func (s *UserServiceServer) Delete(ctx context.Context, req *user.DeleteReq) (*user.DeleteRes, error) {
	err := s.repository.Delete(context.TODO(), req.Id)
	if err != nil {
		return &user.DeleteRes{Status: &user.Status{
			Code: user.Status_FAILED,
			Msg:  "failed",
		}}, err
	}
	return &user.DeleteRes{Status: &user.Status{
		Code: user.Status_SUCCESS,
		Msg:  "success",
	}}, nil
}

func (s *UserServiceServer) List(ctx context.Context, req *user.ListReq) (*user.ListUserResponse, error) {
	cachedUsers, err := s.cache.Get(ctx, "users").Bytes()
	if err != nil {
		users, err := s.repository.FindAll(context.TODO())
		if err != nil {
			return &user.ListUserResponse{
				User: nil,
				Status: &user.Status{
					Code: user.Status_FAILED,
					Msg:  "failed",
				},
			}, err
		}
		err = s.cache.Set(ctx, "users", users, 60*time.Second).Err()
		return &user.ListUserResponse{
			User: users,
			Status: &user.Status{
				Code: user.Status_SUCCESS,
				Msg:  "success",
			},
		}, nil
	}
	var users []*user.User
	err = json.Unmarshal(cachedUsers, &users)

	if err != nil {
		return &user.ListUserResponse{
			User: nil,
			Status: &user.Status{
				Code: user.Status_FAILED,
				Msg:  "failed",
			},
		}, err
	}

	return &user.ListUserResponse{
		User: users,
		Status: &user.Status{
			Code: user.Status_SUCCESS,
			Msg:  "success",
		},
	}, nil
}
func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := UserServiceServer{}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &s)

	log.Println("server start")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
