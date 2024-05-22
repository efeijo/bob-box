package jwt

import (
	"authservice/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: create errors constant

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret = []byte("bobbox key")

var (
	tokenDuration = time.Hour
)

type Validator struct {
	secret []byte
}

type ValidatorConfig struct {
	Secret []byte
}

func NewValidator(config ValidatorConfig) *Validator {
	if config.Secret == nil {
		config.Secret = hmacSampleSecret
	}
	return &Validator{
		secret: config.Secret,
	}
}

// Validate takes in jwt token
func (jv *Validator) Validate(tokenReceived string) (*model.UserClaims, error) {
	token, err := jwt.Parse(tokenReceived, jv.keyFunc)
	if err != nil {
		return nil, err
	}

	return parseAndCheckTokenClaims(token)
}

func (jv *Validator) CreateToken(userID string, exp *time.Duration) (string, error) {
	// TODO: make this a bit prettier
	if exp == nil {
		exp = &tokenDuration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uid": userID,
			"exp": time.Now().Add(*exp).Unix(),
		},
	)

	return token.SignedString(jv.secret)
}

func (jv *Validator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return jv.secret, nil
}

func parseAndCheckTokenClaims(token *jwt.Token) (*model.UserClaims, error) {
	var userID string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id, ok := claims["uid"]
		if !ok {
			return nil, errors.New("no user id found in token claims")
		}

		userID = id.(string)

		exp, err := claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		if !exp.Time.After(time.Now()) {
			return nil, errors.New("token expired")
		}

	} else {
		return nil, errors.New("casting claims failed")
	}
	fmt.Println(userID)
	return &model.UserClaims{UserID: userID}, nil
}
