package ctrl

import "github.com/gin-gonic/gin"

type Controller interface {
	Register(*gin.RouterGroup)
}
