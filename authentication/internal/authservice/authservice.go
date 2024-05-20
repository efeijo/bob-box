package authservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"authservice/internal/authorization"
	"authservice/internal/model"
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
		log.Println("error while getting user from database", err)
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

	fmt.Println(user)

	jwtToken, err := a.validator.CreateToken(user.UserId, time.Hour)
	if err != nil {
		return "", err
	}

	err = a.store.CreateSession(ctx, &model.Session{
		UID:      user.UserId,
		JWTToken: jwtToken,
	})

	// create and respond with jwt token
	return jwtToken, err
}

func (a *Auth) InvalidateToken(ctx context.Context, userID string) error {
	return a.store.DeleteSession(ctx, userID)
}

func (a *Auth) ValidateToken(ctx context.Context, jwtToken string) (bool, error) {
	userID, ok := a.validator.Validate(jwtToken)
	if !ok {
		return false, errors.New("invalid token")
	}
	session, err := a.store.GetSession(ctx, userID)
	if err != nil {
		return false, errors.New("no session found for that token")
	}
	return session.JWTToken == jwtToken, nil
}

func (a *Auth) CreateUser(ctx context.Context, username string, password string) error {
	_, err := a.store.GetUser(ctx, username)
	if err != nil {
		return err
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

func (a *Auth) DeleteUser(ctx context.Context, userID string) error {
	user, err := a.store.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	return a.store.DeleteUser(ctx, user.UserId)
}

func NewAuthService(validator authorization.Validator, store Store) *Auth {
	return &Auth{
		validator: validator,
		store:     store,
	}
}
