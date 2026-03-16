package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("Invalid Authorization header")
	}
	splited := strings.Split(apiKey, " ")
	if len(splited) != 2 || splited[0] != "ApiKey" {
		return "", fmt.Errorf("Invalid Authorization header")
	}
	return splited[1], nil
}
