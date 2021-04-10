package utils

import (
	"errors"
	"net/http"
	"regexp"
)

func GetBearerToken(r http.Request) (*string, error) {
	token := ""
	auth := r.Header.Get("Authorization")
	reg := regexp.MustCompile(`Bearer\s(?P<bearer_token>[^\s]+)`)
	strs := reg.FindStringSubmatch(auth)
	if strs == nil {
		return nil, errors.New("Found no authorization bearer token")
	}
	if len(strs) == 2 {
		token = strs[1]
		return &token, nil
	}
	return nil, errors.New("Found no authorization bearer token")
}
