package auth

import (
	"net/http"

	"go.uber.org/zap"
)

func authorize(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		logger.Error("Missing Authorizartion header")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	claims, err := DecodeJwt(token)
	if err != nil {
		logger.Error("Invalid token", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claims.Issuer != Issuer {
		logger.Error("Invalid token", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
