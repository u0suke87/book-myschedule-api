package main

import (
	"github.com/gin-gonic/gin"
	"github.com/u0suke87/book-myschedule-api/controllers"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.tmpl")
	router.Static("/assets", "./assets")
	router.GET("/", controllers.HomeHandler)
	router.POST("/thanks", controllers.AddCalender)

	router.Run(":8080")
}
