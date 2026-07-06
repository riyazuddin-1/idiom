package main

import (
    "github.com/gin-gonic/gin"
    "idiom-api-services/api/auth"
)

func main() {
    router := gin.Default()
	
    api := router.Group("/api/v1")
    routes.Mount(api)

    router.Run(":8080")
}