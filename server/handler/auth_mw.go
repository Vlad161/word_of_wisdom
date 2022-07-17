package handler

import (
	"net/http"
)

func AuthMW(next http.Handler, tStorage TokenStorage) http.HandlerFunc {
	const bearerPrefix = "Bearer "

	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Header.Get("Authorization")
		if len(v) > len(bearerPrefix) && tStorage.Use(v[len(bearerPrefix):]) == nil {
			next.ServeHTTP(w, req)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
