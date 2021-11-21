package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func (s *service) Challenge(w http.ResponseWriter, r *http.Request) {
	cmd, err := exec.Command("/bin/bash", s.scriptPath+"/scripts/challenge.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	ch := strings.TrimSpace(string(cmd))
	key := base64.StdEncoding.EncodeToString([]byte(ch))
	s.cache.Set(key, Challenge(ch))
	w.Write([]byte(ch))
}

func (s *service) Response(key string) (res string, err error) {
	if cachedChallenge, ok := s.cache.Get(key); ok {
		cmd, err := exec.Command("/bin/bash", s.scriptPath+"/scripts/response.sh", string(cachedChallenge.Challenge)).Output()
		if err != nil {
			fmt.Printf("error %s", err)
		}
		return strings.TrimSpace(string(cmd)), nil
	} else {
		return "", errors.New("challenge not found")
	}
}
