package transport

import (
	"encoding/json"
	"log"
	"net/http"
)

// https://github.com/golang/crypto/tree/master/bcrypt

type ApiError struct {
	err            error `json:"err,omitempty"`
	httpStatusCode int   `json:"http_status_code,omitempty"`
}

func WriteError(w http.ResponseWriter, apiError ApiError) http.HandlerFunc {
	log.Println("Write error handler")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(apiError.httpStatusCode)
		err := json.NewEncoder(w).Encode(
			struct {
				Err error `json:"error"`
			}{
				Err: apiError.err,
			},
		)
		log.Println("error encoding error on WriteError", err)
	}
}

func WriteJson(w http.ResponseWriter, statusCode int, val any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if val != nil {
			err := json.NewEncoder(w).Encode(val)
			if err != nil {
				log.Println("error encode value o WriteJson", err)
			}
		}
		w.Write([]byte{})
	}
}
