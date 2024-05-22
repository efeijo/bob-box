package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"authservice/internal/authservice"
	"authservice/internal/model"
)

type Server struct {
	authService authservice.AuthService
	mux         *chi.Mux
}

func NewServer(service authservice.AuthService) *Server {
	s := &Server{authService: service, mux: chi.NewRouter()}

	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)

	authRouter := s.createAuthSubRouter()

	s.mux.Mount("/auth", authRouter)

	return s
}

func (s *Server) createAuthSubRouter() *chi.Mux {
	subRouter := chi.NewRouter()

	subRouter.Post("/user", s.createUserHandler)
	subRouter.Delete("/user", s.deleteUser)
	subRouter.Get("/users", s.listUsersHandler)

	subRouter.Post("/token", s.getToken)
	subRouter.Delete("/token", s.invalidateToken)
	subRouter.Get("/token/{jwt_token}", s.validateToken)

	return subRouter
}

type UserAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s *Server) getToken(w http.ResponseWriter, req *http.Request) {
	var userAuth UserAuth
	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := s.authService.GetUserToken(req.Context(), userAuth.Username, userAuth.Password)
	if err != nil {
		log.Println("error getting user token", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, http.StatusCreated, struct {
		JwtToken string `json:"jwt_token"`
	}{
		JwtToken: token,
	})

}

func (s *Server) createUserHandler(w http.ResponseWriter, req *http.Request) {
	var userAuth UserAuth
	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.authService.CreateUser(req.Context(), userAuth.Username, userAuth.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("writing response on createUserHandler")
	WriteJson(w, http.StatusCreated, nil)
}

type DeleteUserRequest struct {
	Username string `json:"username,omitempty"`
}

func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	var deleteRequest DeleteUserRequest
	err := json.NewDecoder(req.Body).Decode(&deleteRequest)
	if err != nil {
		WriteError(w, ApiError{err: err, httpStatusCode: http.StatusBadRequest})
	}
	err = s.authService.DeleteUser(req.Context(), deleteRequest.Username)
	if err != nil {
		log.Println("error deleting user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	WriteJson(w, http.StatusOK, []byte{})
}

type InvalidateToken struct {
	Username string `json:"username,omitempty"`
}

func (s *Server) invalidateToken(w http.ResponseWriter, req *http.Request) {
	// TODO: read user from token and deletes it
	var invalidateRequest InvalidateToken
	err := json.NewDecoder(req.Body).Decode(&invalidateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	s.authService.InvalidateToken(req.Context(), invalidateRequest.Username)
}

func (s *Server) listUsersHandler(w http.ResponseWriter, req *http.Request) {
	users, err := s.authService.ListUsers(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if users == nil {
		users = []*model.User{}
	}
	WriteJson(w, http.StatusOK, users)
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Server) validateToken(w http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "jwt_token")
	fmt.Println(token)

	isValid, err := s.authService.ValidateToken(req.Context(), token)
	fmt.Println(isValid, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, http.StatusOK,
		struct {
			IsValid bool
		}{
			IsValid: isValid,
		},
	)
}
