package authorization

import (
	"authservice/internal/model"
	"time"
)

type Validator interface {
	Validate(token string) (*model.UserClaims, error)
	CreateToken(userID string, exp *time.Duration) (string, error)
}
