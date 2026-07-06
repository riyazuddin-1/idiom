package main

import (
	"github.com/gin-gonic/gin",
	"api/auth"
)

func main() {
	router := gin.Default()
	
	api = router.Group("/api")
	Auth(api)

	router.Run(":8080")
}