package http

import (
	"file_storage_service/internal/app/usecase"
	"file_storage_service/internal/domain/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	MailUsercase *usecase.MailUsecase
}

func (h *MailHandler) SendMail(c *gin.Context) {
	var req model.MailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	fmt.Println("Request to send email:", req)
	if err := h.MailUsercase.SendMail(req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Email sent successfully"})
}
