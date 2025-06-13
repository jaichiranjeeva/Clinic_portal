package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"portal/config"
	"portal/routes"
)

func main() {
	config.InitDB()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(r)

	r.Run(":8088")
}
