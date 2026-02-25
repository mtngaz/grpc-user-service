package storage

import (
	"context"
	"fmt"
	"time"

	pb "github.com/mtngaz/grpc-user-service/api"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisStore{client: rdb}
}

func (r *RedisStore) CreateUser(ctx context.Context, user *pb.User) error {

	data, err := proto.Marshal(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("user:%d", user.Id)
	return r.client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (r *RedisStore) GetUser(ctx context.Context, id int64) (*pb.User, error) {
	key := fmt.Sprintf("user:%d", id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user pb.User
	if err := proto.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RedisStore) UpdateUser(ctx context.Context, user *pb.User) error {
	return r.CreateUser(ctx, user)
}

func (r *RedisStore) DeleteUser(ctx context.Context, id int64) (bool, error) {
	key := fmt.Sprintf("user:%d", id)
	result, err := r.client.Del(ctx, key).Result()
	return result > 0, err
}

func (r *RedisStore) GetAllUsers(ctx context.Context) ([]*pb.User, error) {
	keys, err := r.client.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}

	users := make([]*pb.User, 0, len(keys))

	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var user pb.User

		if err := proto.Unmarshal(data, &user); err != nil {
			continue
		}

		users = append(users, &user)
	}

	return users, nil
}
