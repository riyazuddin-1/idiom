import (
	"github.com/gin-gonic/gin",
	"api/auth/v1/handlers"
)

func authRoutes(v1 *gin.RouterGroup) {
	v1.POST("/login", loginHandler)
	v1.POST("/logout", logoutHandler)
	v1.POST("/register", registerHandler)
	v1.PATCH("/password-reset", passwordResetHandler)

	v1.GET("/me", getCurrentUserHandler)
	v1.PUT("/me", updateCurrentUserHandler)
}