package main

import (
	"fmt"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/SiddharthSharechat/CRUDGo/Models"
)

func init() {
	Initializers.LoadEnvVariables()
	Initializers.InitDb()
}

func main() {
	err := Initializers.Db.AutoMigrate(&Models.User{})
	if err != nil {
		fmt.Printf("Migration Failed with %s", err)
		return
	}
}
