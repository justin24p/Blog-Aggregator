package auth

import (
	"errors"
	"net/http"
	"strings"
)

// check if it has an api key in request
func GetAPIKey(headers http.Header) (string,error) {
	val := headers.Get("Authorization")	
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0]  != "Apikey"{
		return "", errors.New("malformed first part auth header")
	}
	return vals[1], nil
}