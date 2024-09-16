package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKeyFromHeaders(header http.Header) (string, error) {
	var authorization string = header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("No authorization found")
	}
	vals_arr := strings.Split(authorization, " ")
	if len(vals_arr) != 2 {
		return "", errors.New("invaid header")
	}
	if vals_arr[0] != "ApiKey" {
		return "", errors.New("invalid firt part of auth header")
	}

	return vals_arr[1], nil
}
