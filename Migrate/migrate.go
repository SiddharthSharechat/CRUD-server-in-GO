package main

import (
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/SiddharthSharechat/CRUDGo/Models"
)

func init() {
	Initializers.LoadEnvVariables()
	Initializers.InitDb()
}

func main() {
	Initializers.Db.AutoMigrate(&Models.User{})
}
