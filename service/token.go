package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenRequestBody struct {
	Key      string
	Response string
}

type TokenResponseBody struct {
	AccessToken string `json:"access_token,omitempty"`
}

func (s *service) Token(w http.ResponseWriter, r *http.Request) {
	tokenShouldBeValid := true
	requestBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var requestBody TokenRequestBody
	err = json.Unmarshal(requestBodyBytes, &requestBody)
	if err != nil {
		tokenShouldBeValid = false
		fmt.Println(err)
	}

	response, err := s.Response(requestBody.Key)
	if err != nil {
		tokenShouldBeValid = false
		fmt.Println(err)
	}
	defer s.cache.Remove(requestBody.Key)

	if response != requestBody.Response {
		tokenShouldBeValid = false
		fmt.Println("not authorised")
		fmt.Println("invalid response", requestBody.Response)
	}

	expire := os.Getenv("TOKEN_EXPIRE")
	if expire == "" {
		expire = "3600"
	}

	expireTime, err := strconv.Atoi(expire)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signingKey := []byte(s.apiKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + int64(expireTime),
	})

	var signedToken string
	if tokenShouldBeValid {
		signedToken, err = token.SignedString(signingKey)
	} else {
		uuidWithHyphen := uuid.New()
		signedToken, err = token.SignedString([]byte(uuidWithHyphen.String()))
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(TokenResponseBody{
		AccessToken: signedToken,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(responseBody))
}
