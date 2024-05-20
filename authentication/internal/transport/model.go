package transport

type SignILoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
