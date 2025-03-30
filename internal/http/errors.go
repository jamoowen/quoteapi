package http

import (
	"net/http"
)

func internalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal server error"
	}
	http.Error(w, message, http.StatusInternalServerError)
}

func badRequestError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad request"
	}
	http.Error(w, message, http.StatusBadRequest)
}

func unauthorizedError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized!"
	}
	http.Error(w, message, http.StatusUnauthorized)
}
