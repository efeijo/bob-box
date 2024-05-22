package redis

import (
	"context"
	"encoding/json"
	"log"

	"authservice/internal/model"
)

// https://github.com/redis/go-redis/blob/master/example/scan-struct/main.go

func (rdb *Redis) GetUser(ctx context.Context, userID string) (*model.User, error) {
	res, err := rdb.db.Get(ctx, UsersKey+userID).Bytes()
	if err != nil {
		return nil, err
	}

	var user *model.User
	return user, json.Unmarshal(res, &user)
}

func (rdb *Redis) CreateOrSetUser(ctx context.Context, user *model.User) error {
	return rdb.db.Set(ctx, UsersKey+user.UserID, user, 0).Err()
}

func (rdb *Redis) ListUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	res, err := rdb.db.Keys(ctx, UsersKey+"*").Result()
	if err != nil {
		log.Println("error while fetching keys from database", err.Error())
		return nil, err
	}
	for _, key := range res {
		userBytes, err := rdb.db.Get(ctx, key).Bytes()
		if err != nil {
			log.Println("error while fetching user from database", err.Error())
			return nil, err
		}
		var user *model.User
		err = json.Unmarshal(userBytes, &user)
		if err != nil {
			log.Println("error while unmarshalling user from database", err.Error())
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (rdb *Redis) DeleteUser(ctx context.Context, userID string) error {
	return rdb.db.Del(ctx, UsersKey+userID).Err()
}
