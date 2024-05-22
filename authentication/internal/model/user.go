package model

import (
	"encoding/json"
)

type User struct {
	Username       string `json:"username,omitempty"`
	HashedPassword string `json:"password,omitempty"`
	LoggedIn       bool   `json:"logged_in,omitempty"`
	UserID         string `json:"user_id,omitempty"`
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

type UserClaims struct {
	UserID string `json:"user_id"`
}
