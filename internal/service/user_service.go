package service

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/mtngaz/grpc-user-service/api"
	"github.com/mtngaz/grpc-user-service/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	store *storage.RedisStore
}

func NewUserService(store *storage.RedisStore) * UserService {
	return &UserService{store: store}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &pb.User{
		Id: time.Now().UnixNano(),
		Name: req.Name,
		Email: req.Email,
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		log.Printf("[CreatorUser] error: %v", err)
	}

	log.Printf("[CreatorUser] id = %d name = %s", user.Id, user.Name)

	return &pb.CreateUserResponse{User: user}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := &pb.User{
		Id: req.Id,
		Name: req.Name,
		Email: req.Email,
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		log.Printf("[UpdateUser] error: %v", err)
	}

	log.Printf("[UpdateUser] id = %d name = %s", user.Id, user.Name)
	return &pb.UpdateUserResponse{User: user}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	deleted, err := s.store.DeleteUser(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
    }

    return &pb.DeleteUserResponse{
        Success: deleted,
        Message: map[bool]string{
			true:  fmt.Sprintf("User %d deleted", req.Id),
            false: fmt.Sprintf("User %d not found", req.Id), 
		}[deleted],
    }, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.store.GetAllUsers(ctx)
	if err != nil {
		log.Printf("[GetAllUsers] error: %v", err)
		return nil, err
	}

	log.Printf("[GetAllUsers] count=%d", len(users))
	return &pb.GetAllUsersResponse{Users: users}, nil
}