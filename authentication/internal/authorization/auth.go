package authorization

import (
	"time"
)

type Validator interface {
	Validate(token string) (string, bool)
	CreateToken(userID string, nbf time.Duration) (string, error)
}
