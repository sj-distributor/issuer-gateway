package handler

import "github.com/gin-gonic/gin"

func ApiPrefix(gin *gin.Engine) *gin.RouterGroup {
	return gin.Group("api")
}
