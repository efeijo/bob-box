package main

import (
	"log"

	"authservice/internal/authorization/jwt_auth"
	"authservice/internal/authservice"
	"authservice/internal/store"
	"authservice/internal/transport"
)

func main() {

	client := store.NewRedisClient()

	jwtValidator := jwt_auth.NewJwtValidator(
		jwt_auth.JwtValidatorConfig{
			Secret: []byte("some-secret"),
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
