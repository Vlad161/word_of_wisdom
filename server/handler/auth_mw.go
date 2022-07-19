package handler

import (
	"context"
	"net/http"
)

const bearerPrefix = "Bearer "

func AuthMW(next http.Handler, tStorage TokenStorage, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if token := extractHeaderBearer(req.Header); len(token) > 0 {
			if verifyBearer(req.Context(), token, tStorage, jwt) {
				next.ServeHTTP(w, req)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

func verifyBearer(ctx context.Context, token string, tStorage TokenStorage, jwt JWTService) bool {
	payload, err := jwt.Verify(token)
	if err != nil {
		return false
	}

	isVerify, okVerify := payload[keyIsVerify].(bool)
	tc, okTc := payload[keyToken].(string)
	if okVerify && okTc {
		return isVerify && tStorage.Use(ctx, tc) == nil
	}

	return false
}

func extractHeaderBearer(h http.Header) string {
	v := h.Get("Authorization")
	if len(v) > len(bearerPrefix) {
		return v[len(bearerPrefix):]
	}
	return ""
}
