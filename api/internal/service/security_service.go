package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/model"
)

type SecurityService interface {
	VerifyTurnstile(token string, remoteIP string) error
	SanitizeContent(content string) string
}

type securityService struct {
	turnstileSecret string
	sanitizer       *bluemonday.Policy
	httpClient      *http.Client
}

func NewSecurityService(secret string) SecurityService {
	return &securityService{
		turnstileSecret: secret,
		// UGCPolicy mengizinkan formatting dasar tapi menghapus script/iframe berbahaya
		sanitizer: bluemonday.UGCPolicy(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Prevent hanging requests
		},
	}
}

func (s *securityService) VerifyTurnstile(token string, remoteIP string) error {
	resp, err := s.httpClient.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify",
		url.Values{
			"secret":   {s.turnstileSecret},
			"response": {token},
			"remoteip": {remoteIP},
		})
	if err != nil {
		return errors.New("failed to connect to turnstile provider")
	}
	defer resp.Body.Close()

	var result model.TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return errors.New("failed to parse turnstile response")
	}

	if !result.Success {
		return errors.New("invalid captcha token")
	}

	return nil
}

func (s *securityService) SanitizeContent(content string) string {
	return s.sanitizer.Sanitize(content)
}