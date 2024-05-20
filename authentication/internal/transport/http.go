package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"authservice/internal/authservice"
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
	subRouter.Get("/token", s.validateToken)

	return subRouter
}

func (s *Server) getToken(w http.ResponseWriter, req *http.Request) {
	token, err := s.authService.GetUserToken(req.Context(), "Emanuel", "1234567")
	if err != nil {
		ApiError(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(
		struct {
			JwtToken string `json:"jwt_token"`
		}{
			JwtToken: token,
		},
	)

	if err != nil {
		ApiError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) createUserHandler(w http.ResponseWriter, req *http.Request) {
	err := s.authService.CreateUser(req.Context(), "Emanuel", "1234567")
	if err != nil {
		ApiError(w, err, http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	err := s.authService.DeleteUser(req.Context(), "Emanuel")
	if err != nil {
		ApiError(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Server) invalidateToken(w http.ResponseWriter, req *http.Request) {
	err := s.authService.InvalidateToken(req.Context(), "Emanuel")
	if err != nil {
		ApiError(w, err, http.StatusNotFound)
		return
	}
}

func (s *Server) listUsersHandler(w http.ResponseWriter, req *http.Request) {
	users, err := s.authService.ListUsers(req.Context())
	if err != nil {
		ApiError(w, err, http.StatusInternalServerError)
		return
	}

	if users == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		ApiError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Server) validateToken(w http.ResponseWriter, req *http.Request) {
	isValid, err := s.authService.ValidateToken(req.Context(), "token here")
	if err != nil {
		return
	}
	fmt.Println(isValid)
}
