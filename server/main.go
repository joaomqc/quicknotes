package main

import (
	"fmt"
	"os"
	"quicknotes/ctrl"
	"quicknotes/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const DEFAULT_PORT = "8000"

//	@title			quicknotes API
//	@version		1.0
//	@description	Note management app
//	@BasePath		/api

var ctrls = []ctrl.Controller{
	ctrl.NotesController{},
}

func main() {
	r := gin.Default()

	apiGroup := r.Group("/api")
	for _, controller := range ctrls {
		controller.Register(apiGroup)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("GOPHERSS_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		panic("[ERROR] failed to start server: " + err.Error())
	}
}
