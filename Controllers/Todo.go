package Controllers

import (
	"fmt"
	"net/http"

	"go-todo-app/Config"
	"go-todo-app/Models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	var todos []Models.Todo
	db := Config.ConnectToDB()
	defer db.Close()
	row, err := db.Query("SELECT * FROM todo")
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	for row.Next() {
		var todo Models.Todo
		if err := row.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			fmt.Fprint(c.Writer, err)
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}

func CreateATodo(c *gin.Context) {
	var todo Models.Todo
	c.BindJSON(&todo)
	db := Config.ConnectToDB()
	defer db.Close()
	_, err := db.Query("insert into todo(ID, Title, Description) values(?,?,?)", todo.ID, todo.Title, todo.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, "Todo created Successfully.....")
}

func GetATodo(c *gin.Context) {
	id := c.Params.ByName("id")
	var todo Models.Todo
	db := Config.ConnectToDB()
	defer db.Close()
	row, err := db.Query("SELECT * FROM todo where ID=?", id)
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	for row.Next() {
		if err := row.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			fmt.Fprint(c.Writer, err)
			return
		}
	}
	// c.BindJSON(row)
	c.JSON(http.StatusOK, todo)
}

func UpdateATodo(c *gin.Context) {
	id := c.Params.ByName("id")
	var todo Models.Todo
	c.BindJSON(&todo)
	db := Config.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("update todo set Title=?, Description=? where ID=?", todo.Title, todo.Description, id)
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	c.JSON(http.StatusOK, "Updated Successfully.......")
}

func DeleteATodo(c *gin.Context) {
	id := c.Params.ByName("id")
	db := Config.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("DELETE from todo where ID=?", id)
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	c.JSON(http.StatusOK, "Record deleted Succesfully.......")
}
