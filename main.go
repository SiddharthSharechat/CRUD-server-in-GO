package main

import (
	"fmt"
	"github.com/SiddharthSharechat/CRUDGo/Controllers"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	Initializers.LoadEnvVariables()
	Initializers.InitDb()
	Initializers.InitCache()
}

func main() {
	r := gin.Default()

	r.POST("/user", Controllers.UsersCreate)
	r.GET("/user/:id", Controllers.UsersGet)
	r.PUT("/user/:id", Controllers.UsersUpdate)
	r.DELETE("/user/:id", Controllers.UsersDelete)
	r.GET("/user", Controllers.UsersGetPaginated)

	err := r.Run()
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
		return
	}
}
