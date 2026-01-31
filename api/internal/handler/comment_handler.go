package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/model"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/service"
)

type CommentHandler struct {
	securityService service.SecurityService
}

func NewCommentHandler(sec service.SecurityService) *CommentHandler {
	return &CommentHandler{
		securityService: sec,
	}
}

func (h *CommentHandler) PostComment(c *gin.Context) {
	var req model.CommentRequest

	// 1. Structural Validation (Binding & Validator)
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return general error untuk user, log detail untuk developer
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	// 2. Turnstile Verification (Server-side)
	clientIP := c.ClientIP()
	if err := h.securityService.VerifyTurnstile(req.TurnstileToken, clientIP); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bot verification failed"})
		return
	}

	// 3. XSS Sanitization
	// Bersihkan content sebelum masuk ke database atau logic selanjutnya
	cleanContent := h.securityService.SanitizeContent(req.Content)

	// --- Di sini logic simpan ke database (Repository Layer) ---
	// commentRepo.Save(req.Username, cleanContent)
	
	// Simulasi sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment submitted successfully",
		"data": gin.H{
			"username": req.Username,
			"content":  cleanContent, // Return versi yang sudah dibersihkan
		},
	})
}