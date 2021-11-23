package service_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/golang-jwt/jwt/v4"
	"github.com/meinto/public-api-auth-provider-service/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("token endpoint", func() {

	var s service.Service
	var w *httptest.ResponseRecorder
	correctResponseHash := sha256.New()
	correctResponseHash.Write([]byte("the-response-for-the-challenge the-challenge"))
	correctResponseString := fmt.Sprintf("%x", correctResponseHash.Sum(nil))

	BeforeEach(func() {
		scriptPath, _ := filepath.Abs("../")
		s = service.NewService(service.ServiceOptions{
			ApiKey:     "test-key",
			ScriptPath: scriptPath,
		})

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/challenge", nil)
		w = httptest.NewRecorder()

		s.Challenge(w, req)
	})

	It("returns a valid access token", func() {
		res := w.Result()
		bodyBytes, _ := io.ReadAll(res.Body)
		key := base64.StdEncoding.EncodeToString(bodyBytes)

		postBody, _ := json.Marshal(service.TokenRequestBody{
			Key:      key,
			Response: correctResponseString,
		})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/token", bytes.NewReader(postBody))
		w = httptest.NewRecorder()

		s.Token(w, req)

		res = w.Result()
		bodyBytes, _ = io.ReadAll(res.Body)

		var responseBody service.TokenResponseBody
		json.Unmarshal(bodyBytes, &responseBody)

		token, err := jwt.Parse(responseBody.AccessToken, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-key"), nil
		})

		Expect(res.StatusCode).To(Equal(http.StatusOK))
		Expect(err).To(BeNil())
		Expect(token.Valid).To(BeTrue())
	})

	It("rejects if the response for the challenge is wrong", func() {
		res := w.Result()
		bodyBytes, _ := io.ReadAll(res.Body)
		key := base64.StdEncoding.EncodeToString(bodyBytes)

		postBody, _ := json.Marshal(service.TokenRequestBody{
			Key:      key,
			Response: "wrong response",
		})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/token", bytes.NewReader(postBody))
		w = httptest.NewRecorder()
		s.Token(w, req)

		// second call
		w = httptest.NewRecorder()
		s.Token(w, req)

		res = w.Result()

		Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("rejects if the key for the challenge is wrong", func() {
		postBody, _ := json.Marshal(service.TokenRequestBody{
			Key:      "wrong-key",
			Response: correctResponseString,
		})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/token", bytes.NewReader(postBody))
		w = httptest.NewRecorder()
		s.Token(w, req)

		// second call
		w = httptest.NewRecorder()
		s.Token(w, req)

		res := w.Result()

		Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("rejects if the challenge is already solved", func() {
		res := w.Result()
		bodyBytes, _ := io.ReadAll(res.Body)
		key := base64.StdEncoding.EncodeToString(bodyBytes)

		postBody, _ := json.Marshal(service.TokenRequestBody{
			Key:      key,
			Response: correctResponseString,
		})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/token", bytes.NewReader(postBody))
		w = httptest.NewRecorder()
		s.Token(w, req)

		// second call
		w = httptest.NewRecorder()
		s.Token(w, req)

		res = w.Result()

		Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
