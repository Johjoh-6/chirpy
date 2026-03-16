package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
	if bearerToken == "" {
		return "", errors.New("invalid bearer token")
	}

	return bearerToken, nil
}
