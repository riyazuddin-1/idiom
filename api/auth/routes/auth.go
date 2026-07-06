package auth

import (
	"github.com/gin-gonic/gin"
	"idiom-api-services/api/auth/handlers"
)

func Mount(router *gin.RouterGroup) {
	router.POST("/login", handlers.LoginHandler)
	router.POST("/logout", handlers.LogoutHandler)
	router.POST("/register", handlers.RegisterHandler)
	router.PATCH("/password-reset", handlers.PasswordResetHandler)

	router.GET("/me", handlers.GetCurrentUserHandler)
	router.PUT("/me", handlers.UpdateCurrentUserHandler)
}
