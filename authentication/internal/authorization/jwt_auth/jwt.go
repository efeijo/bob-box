package jwt_auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret []byte = []byte("bobbox key")

var (
	tokenDuration = time.Hour
)

type JwtValidator struct {
	secret []byte
}

type JwtValidatorConfig struct {
	Secret []byte
}

func NewJwtValidator(config JwtValidatorConfig) *JwtValidator {
	if config.Secret == nil {
		config.Secret = hmacSampleSecret
	}
	return &JwtValidator{
		secret: config.Secret,
	}
}

func (jv *JwtValidator) Validate(tokenReceived string) (string, bool) {
	token, err := jwt.Parse(tokenReceived, func(token *jwt.Token) (interface{}, error) {
		// validates alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jv.secret, nil
	})
	if err != nil {
		fmt.Println("error parsing token", err)
		return "", false
	}

	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		uid = claims["uid"].(string)
		duration := claims["exp"].(float64)
		if int64(duration) < time.Now().Unix() {
			fmt.Println("token expired")
			return "", false
		}
	} else {
		fmt.Println("casting claims failed")
		return "", false
	}

	return uid, true
}

func (jv *JwtValidator) CreateToken(userID string, nbf time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(nbf).Unix(),
	})

	tokenString, err := token.SignedString(jv.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
