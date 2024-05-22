package main

import (
	"log"

	"authservice/internal/authorization/jwt"
	"authservice/internal/authservice"
	"authservice/internal/store/redis"
	"authservice/internal/transport"
)

func main() {
	secret := []byte("some-secret")

	client := redis.NewClient()

	jwtValidator := jwt.NewValidator(
		jwt.ValidatorConfig{
			Secret: secret,
		},
	)

	authService := authservice.NewAuthService(
		jwtValidator,
		client,
	)

	if err := transport.NewServer(authService).ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
