package service_test

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/meinto/public-api-auth-provider-service/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("challange endpoint", func() {
	It("returns a challange & persists it in the inmemory cache", func() {
		scriptPath, _ := filepath.Abs("../")
		s := service.NewService(service.ServiceOptions{
			ApiKey:     "test-key",
			ScriptPath: scriptPath,
		})

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/challenge", nil)
		w := httptest.NewRecorder()

		s.Challenge(w, req)

		res := w.Result()
		bodyBytes, _ := io.ReadAll(res.Body)
		challange := string(bodyBytes)
		key := base64.StdEncoding.EncodeToString(bodyBytes)

		cacheEntry, ok := s.Cache().Get(key)

		Expect(challange).To(Equal("the-challenge"))
		Expect(cacheEntry.Challenge).To(Equal(service.Challenge(challange)))
		Expect(ok).To(BeTrue())
	})
})

var _ = Describe("response function", func() {

	correctResponseHash := sha256.New()
	correctResponseHash.Write([]byte("the-response-for-the-challenge the-challenge"))
	correctResponseString := fmt.Sprintf("%x", correctResponseHash.Sum(nil))

	It("returns the response corresponding to the challenge", func() {
		scriptPath, _ := filepath.Abs("../")
		s := service.NewService(service.ServiceOptions{
			ApiKey:     "test-key",
			ScriptPath: scriptPath,
		})

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/challenge", nil)
		w := httptest.NewRecorder()

		s.Challenge(w, req)

		res := w.Result()
		bodyBytes, _ := io.ReadAll(res.Body)
		key := base64.StdEncoding.EncodeToString(bodyBytes)

		result, err := s.Response(key)

		Expect(err).To(BeNil())
		Expect(result).To(Equal(correctResponseString))
	})
})
