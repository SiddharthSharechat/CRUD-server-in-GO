package main

import (
	"github.com/SiddharthSharechat/CRUDGo/Controllers"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	Initializers.LoadEnvVariables()
	Initializers.InitDb()
}

func main() {
	r := gin.Default()

	r.POST("/users", Controllers.UsersCreate)
	r.GET("/users/:id", Controllers.UsersGet)
	r.PUT("users/:id", Controllers.UsersUpdate)
	r.DELETE("users/:id", Controllers.UsersDelete)
	r.GET("users", Controllers.UsersGetPaginated)

	r.Run()
}
