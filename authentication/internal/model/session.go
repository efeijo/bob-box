package model

type Session struct {
	UID      string `json:"uid"`
	JWTToken string `json:"token"`
}
