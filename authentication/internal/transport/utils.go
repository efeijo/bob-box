package transport

import (
	"net/http"
)

// https://github.com/golang/crypto/tree/master/bcrypt

func ApiError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
