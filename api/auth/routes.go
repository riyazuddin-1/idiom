func Auth(api *gin.RouterGroup) {
	v1 = api.Group("/v1")
	v1Routes(v1)
}