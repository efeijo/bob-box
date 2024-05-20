package store

import (
	"context"
	"encoding/json"

	"authservice/internal/model"
)

func (rdb *Redis) CreateSession(ctx context.Context, session *model.Session) error {
	return rdb.db.Set(ctx, SessionKey+session.UID, session, 0).Err()
}

func (rdb *Redis) GetSession(ctx context.Context, uid string) (*model.Session, error) {
	var session *model.Session

	bytes, err := rdb.db.Get(ctx, SessionKey+uid).Bytes()
	if err != nil {
		return nil, err
	}

	return session, json.Unmarshal(bytes, session)
}

func (rdb *Redis) DeleteSession(ctx context.Context, uid string) error {
	return rdb.db.Del(ctx, SessionKey+uid).Err()
}
