package main

import (
	"path/filepath"

	"github.com/meinto/public-api-auth-provider-service/service"
)

func main() {
	scriptPath, _ := filepath.Abs("./")
	service.NewService(service.ServiceOptions{
		ApiKey:     "api-key",
		ScriptPath: scriptPath,
	}).RunAndServe()
}
