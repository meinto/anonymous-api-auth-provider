package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	Response(key string) (res string, err error)
	Token(w http.ResponseWriter, r *http.Request)
	Challenge(w http.ResponseWriter, r *http.Request)
	RunAndServe()
	Cache() *Cache
}

type ServiceOptions struct {
	ApiKey     string
	ScriptPath string
}

type service struct {
	apiKey     string
	scriptPath string
	cache      *Cache
	handler    *mux.Router
}

func NewService(options ServiceOptions) Service {
	return &service{
		apiKey:     options.ApiKey,
		scriptPath: options.ScriptPath,
		cache:      NewCache(),
		handler:    mux.NewRouter().StrictSlash(true),
	}
}

func (s *service) Cache() *Cache {
	return s.cache
}

func (s *service) RunAndServe() {
	s.handler.HandleFunc("/challenge", s.Challenge)
	s.handler.HandleFunc("/token", s.Token).Methods("POST")

	fmt.Println("listen on port :8080")
	if err := http.ListenAndServe(":8080", s.handler); err != nil {
		log.Fatal("cloud not start server")
	}
}
