package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/meinto/public-api-auth-provider-service/service"
)

func main() {
	scriptPath, _ := filepath.Abs("./")
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		log.Fatal("please set your api key via the API_KEY environment variable")
	}

	service.NewService(service.ServiceOptions{
		ApiKey:     apiKey,
		ScriptPath: scriptPath,
	}).RunAndServe()
}
