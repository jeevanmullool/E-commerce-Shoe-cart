package main

import (
	"redkart/initializers"
	"redkart/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	routes.UserRoutes(r)
	routes.AdminRoutes(r)
	r.Run()

}
