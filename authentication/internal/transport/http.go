package transport

import (
	"encoding/json"
	"log"
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
		log.Println("error getting user token", err)
		WriteError(w, ApiError{err: err, httpStatusCode: http.StatusInternalServerError})
		return
	}

	WriteJson(w, http.StatusCreated, struct {
		JwtToken string `json:"jwt_token"`
	}{
		JwtToken: token,
	})

}

func (s *Server) createUserHandler(w http.ResponseWriter, req *http.Request) {
	err := s.authService.CreateUser(req.Context(), "Emanuel", "1234567")
	if err != nil {
		log.Println("error creating user password", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("writing response on createUserHandler")
	WriteJson(w, http.StatusCreated, nil)
}

func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	err := s.authService.DeleteUser(req.Context(), "Emanuel")
	if err != nil {
		log.Println("error deleting user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	WriteJson(w, http.StatusOK, []byte{})
}

func (s *Server) invalidateToken(w http.ResponseWriter, req *http.Request) {
	s.authService.InvalidateToken(req.Context(), "Emanuel")
}

func (s *Server) listUsersHandler(w http.ResponseWriter, req *http.Request) {
	users, err := s.authService.ListUsers(req.Context())
	if err != nil {
		WriteError(w, ApiError{httpStatusCode: http.StatusInternalServerError, err: err})
	}

	if users == nil {
		WriteJson(w, http.StatusOK, []string{})
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		WriteError(w, ApiError{err: err, httpStatusCode: http.StatusInternalServerError})
	}

	WriteJson(w, http.StatusOK, users)
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Server) validateToken(w http.ResponseWriter, req *http.Request) {
	isValid, err := s.authService.ValidateToken(req.Context(), "token here")
	if err != nil {
		WriteError(w, ApiError{err: err, httpStatusCode: http.StatusInternalServerError})
	}

	WriteJson(w, http.StatusOK,
		struct {
			IsValid bool
		}{
			IsValid: isValid,
		},
	)
}
