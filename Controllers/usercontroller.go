package Controllers

import (
	"go-todo-app/Config"
	models "go-todo-app/Models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	db := Config.ConnectToDB()
	defer db.Close()
	_, err := db.Query("insert into users(Name, Username, Email, Password) values(?,?,?,?)", user.Name, user.Username, user.Email, user.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, "Success........")
}
