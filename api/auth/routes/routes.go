package routes

import (
	"github.com/gin-gonic/gin",
	"idiom-api-services/api/auth/routes/auth",
	"idiom-api-services/api/auth/routes/session",
)

func Mount(router *gin.RouterGroup) {
	auth.Mount(router)
	session.Mount(router)
}
