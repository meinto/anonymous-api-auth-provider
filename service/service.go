package service

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	s.handler.HandleFunc("/challenge", s.Challenge).Methods("GET")
	s.handler.HandleFunc("/token", s.Token).Methods("POST")
	s.handler.HandleFunc("/health", s.Health).Methods("GET")
	s.handler.HandleFunc("/", s.Root)

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	fmt.Printf("listen on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), s.handler); err != nil {
		log.Fatal("cloud not start server")
	}
}

func (s *service) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *service) Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("i'm alive!"))
}
