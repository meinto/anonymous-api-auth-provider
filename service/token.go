package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenRequestBody struct {
	Key      string
	Response string
}

type TokenResponseBody struct {
	AccessToken string `json:"access_token,omitempty"`
}

func (s *service) Token(w http.ResponseWriter, r *http.Request) {
	requestBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var requestBody TokenRequestBody
	err = json.Unmarshal(requestBodyBytes, &requestBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := s.Response(requestBody.Key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer s.cache.Remove(requestBody.Key)

	if response != requestBody.Response {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	signingKey := []byte(s.apiKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + 15000,
	})

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(TokenResponseBody{
		AccessToken: signedToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(responseBody))
}
