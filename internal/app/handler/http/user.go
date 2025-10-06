package http

import (
	"file_storage_service/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUC *usecase.UserUseCase
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.UserUC.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}
