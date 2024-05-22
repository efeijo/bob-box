package authservice

import (
	"context"
	"errors"

	"authservice/internal/authorization"
	"authservice/internal/model"

	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	// token

	GetUserToken(ctx context.Context, username string, password string) (jwtToken string, err error)
	CreateUser(ctx context.Context, username string, password string) error
	InvalidateToken(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, jwtToken string) (bool, error)

	// users

	ListUsers(ctx context.Context) ([]*model.User, error)
	DeleteUser(ctx context.Context, username string) error
}

type Store interface {
	// Sessions

	CreateSession(ctx context.Context, session *model.Session) error
	GetSession(ctx context.Context, uid string) (*model.Session, error)
	DeleteSession(ctx context.Context, uid string) error

	// users

	GetUser(ctx context.Context, uid string) (*model.User, error)
	CreateOrSetUser(ctx context.Context, user *model.User) error
	ListUsers(ctx context.Context) ([]*model.User, error)
	DeleteUser(ctx context.Context, uid string) error
}

type Auth struct {
	validator authorization.Validator
	store     Store
}

func (a *Auth) GetUserToken(ctx context.Context, username string, password string) (string, error) {
	user, err := a.store.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	passwordsMatch := compareHashAndPassword(user.HashedPassword, password)
	if !passwordsMatch {
		return "", errors.New("invalid password given")
	}

	user.LoggedIn = true

	err = a.store.CreateOrSetUser(ctx, user)
	if err != nil {
		return "", err
	}

	jwtToken, err := a.validator.CreateToken(user.Username, nil)
	if err != nil {
		return "", err
	}

	err = a.store.CreateSession(ctx, &model.Session{
		UID:      user.Username,
		JWTToken: jwtToken,
	})

	// create and respond with jwt token
	return jwtToken, err
}

func (a *Auth) InvalidateToken(ctx context.Context, username string) error {
	return a.store.DeleteSession(ctx, username)
}

func (a *Auth) ValidateToken(ctx context.Context, jwtToken string) (bool, error) {
	claimsFromToken, err := a.validator.Validate(jwtToken)
	if err != nil {
		return false, err
	}
	session, err := a.store.GetSession(ctx, claimsFromToken.UserID)
	if err != nil {
		return false, err
	}
	return session.JWTToken == jwtToken, nil
}

func (a *Auth) CreateUser(ctx context.Context, username string, password string) error {
	user, err := a.store.GetUser(ctx, username)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if user != nil {
		return errors.New("user already exists")
	}

	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	return a.store.CreateOrSetUser(ctx, &model.User{
		Username:       username,
		HashedPassword: encryptedPassword,
		LoggedIn:       false,
	})
}

func (a *Auth) ListUsers(ctx context.Context) ([]*model.User, error) {
	return a.store.ListUsers(ctx)
}

func (a *Auth) DeleteUser(ctx context.Context, username string) error {
	user, err := a.store.GetUser(ctx, username)
	if err != nil {
		return err
	}
	return a.store.DeleteUser(ctx, user.Username)
}

func NewAuthService(validator authorization.Validator, store Store) *Auth {
	return &Auth{
		validator: validator,
		store:     store,
	}
}
