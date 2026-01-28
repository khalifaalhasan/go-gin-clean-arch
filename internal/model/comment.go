package model 


type CommentRequest struct {
	TurnstileToken string `json:"turnstile_token" binding:"required"`
	Username       string `json:"username" binding:"required,min=3,max=30,alphanum"`
	Content        string `json:"content" binding:"required,min=5,max=1000"` // Pastikan json tag-nya beda
}

type TurnstileResponse struct {
	Success bool `json:"success"`
	ChallengeTs string `json:"challenge_tss"`
	Hostname string  `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
}