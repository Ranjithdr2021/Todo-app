package main

import (
	"go-todo-app/Config"
	routes "go-todo-app/Routes"

	_ "github.com/go-sql-driver/mysql"
)

var err error

func main() {

	Config.NewTable()

	r := routes.SetupRouter()

	r.Run("localhost:6000")
}
