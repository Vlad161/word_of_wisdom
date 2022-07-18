package handler

import (
	"net/http"
)

const bearerPrefix = "Bearer "

func AuthMW(next http.Handler, tStorage TokenStorage, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if token := extractHeaderBearer(req.Header); len(token) > 0 {
			payload, err := jwt.Verify(token)
			if err == nil && payload[keyIsVerify].(bool) && tStorage.Use(req.Context(), payload[keyToken].(string)) == nil {
				next.ServeHTTP(w, req)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

func extractHeaderBearer(h http.Header) string {
	v := h.Get("Authorization")
	if len(v) > len(bearerPrefix) {
		return v[len(bearerPrefix):]
	}
	return ""
}
