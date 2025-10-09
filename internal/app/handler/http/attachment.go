package http

import (
	"file_storage_service/internal/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	AttachmentUC *usecase.AttachmentUseCase
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	src, _ := file.Open()
	defer src.Close()

	ctx := c.Request.Context()
	key, err := h.AttachmentUC.UploadFile(ctx, src, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"fileKey": key})
}

func (h *AttachmentHandler) Download(c *gin.Context) {
	key := c.Param("key")
	data, err := h.AttachmentUC.DownloadFile(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (h *AttachmentHandler) Delete(c *gin.Context) {
	key := c.Param("key")
	err := h.AttachmentUC.DeleteFile(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
