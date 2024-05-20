package jwt_auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJwtValidator_CreateToken(t *testing.T) {
	jv := &JwtValidator{secret: hmacSampleSecret}

	token, err := jv.CreateToken("emanuel", time.Hour)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(token)

	// parse token
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// validates alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jv.secret, nil
	})

	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		s := claims["uid"].(string)
		if s != "emanuel" {
			t.Error("token claims are wrong")
		}
		parseDuration, ok := claims["exp"].(float64)
		if !ok {
			t.Error("exp claims are wrong")
		}

		if int64(parseDuration) != time.Now().Unix()+3600 {
			t.Error(errors.New("durations"))
		}

	}

	uid, b := jv.Validate(token)
	fmt.Println(uid, b)
	if !b {
		t.Error("token validation failed")
	}
	if uid != "emanuel" {
		t.Error("token claims are wrong")
	}

}
