package redis

import (
	"context"
	"encoding/json"
	"strings"

	"authservice/internal/model"
)

func (rdb *Redis) CreateSession(ctx context.Context, session *model.Session) error {
	return rdb.db.Set(
		ctx,
		SessionKey+session.UID,
		session,
		0,
	).Err()
}

func (rdb *Redis) GetSession(ctx context.Context, userID string) (*model.Session, error) {
	var session *model.Session

	bytes, err := rdb.db.Get(ctx, SessionKey+userID).Bytes()
	if err != nil {
		return nil, err
	}

	return session, json.Unmarshal(bytes, session)
}

func (rdb *Redis) DeleteSession(ctx context.Context, userID string) error {
	return rdb.db.Del(
		ctx,
		SessionKey+userID,
	).Err()
}

func (rdb *Redis) ListSessions(ctx context.Context, userID string) ([]*model.Session, error) {
	ks, err := rdb.db.Keys(
		ctx,
		SessionKey+"*",
	).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*model.Session, 0, len(ks))
	for _, key := range ks {
		userID, _ := strings.CutPrefix(key, SessionKey)

		session, err := rdb.GetSession(ctx, userID)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)

	}

	return sessions, nil
}
