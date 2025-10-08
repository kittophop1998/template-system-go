package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all API routes and their handlers.
func SetupRoutes(router *gin.Engine, handlers *Handlers) {
	registerHealthCheck(router)
	registerAPIV1Routes(router, handlers)
}

// registerHealthCheck handles the health check endpoint.
func registerHealthCheck(router *gin.Engine) {
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// registerAPIV1Routes registers all version 1 API routes.
func registerAPIV1Routes(router *gin.Engine, handlers *Handlers) {
	apiV1 := router.Group("/api/v1")

	registerUserRoutes(apiV1, handlers)
	registerFileRoutes(apiV1, handlers)
	registerMailRoutes(apiV1, handlers)
}

// registerUserRoutes handles user-related endpoints.
func registerUserRoutes(rg *gin.RouterGroup, handlers *Handlers) {
	users := rg.Group("/users")
	users.GET("", handlers.User.GetUsers)
}

// registerFileRoutes handles file upload endpoints.
func registerFileRoutes(rg *gin.RouterGroup, handlers *Handlers) {
	files := rg.Group("/files")
	files.POST("/upload", handlers.Attachment.Upload)
}

// registerMailRoutes handles mail-related endpoints.
func registerMailRoutes(rg *gin.RouterGroup, handlers *Handlers) {
	mails := rg.Group("/mails")
	mails.POST("/send", handlers.Mail.SendMail)
}
