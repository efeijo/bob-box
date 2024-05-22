package model

import "encoding/json"

type Session struct {
	UID      string `json:"uid"`
	JWTToken string `json:"token"`
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Session) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}
