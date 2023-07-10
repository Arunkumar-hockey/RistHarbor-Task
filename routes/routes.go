package routes

import (
	"fmt"
	"test/controller"

	"github.com/gin-gonic/gin"
)


func StartApplication() {
	fmt.Println("Application Started...")
	var router = gin.Default()
	router.POST("api/v1/create", controller.SaveBlocks)
	router.GET("api/v1/query", controller.QueryBlocks)

	router.Run(":8000")
}