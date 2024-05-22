package jwt

import (
	"fmt"
	"testing"
)

func TestJwtValidator_CreateToken(t *testing.T) {
	someID := "emanuel"
	jv := &Validator{secret: hmacSampleSecret}

	token, err := jv.CreateToken(someID, nil)
	if err != nil {
		t.Error(err)
	}

	claims, err := jv.Validate(token)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(claims, err)
	if claims.UserID != someID {
		t.Error("token claims are wrong")
	}

}
