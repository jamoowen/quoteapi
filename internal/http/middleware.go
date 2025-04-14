package http

import "net/http"

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "X-API-KEY header missing", http.StatusForbidden)
			return
		}
		// check db for user
		authorized, err := h.authService.AuthenticateApiKey(apiKey, r.Context())
		if err != nil {
			h.logger.Println("Failed to authenticate api key: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if authorized == false {
			http.Error(w, "Invalid api key", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
