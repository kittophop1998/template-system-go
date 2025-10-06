package http

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handlers *Handlers) {
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	apiV1 := router.Group("/api/v1")
	{
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.GET("", handlers.User.GetUsers)
		}

		files := apiV1.Group("/files")
		{
			files.POST("/upload", handlers.Attachment.Upload)
		}
	}
}
