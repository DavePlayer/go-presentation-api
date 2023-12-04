package handlers

import (
	"apitest/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context, db *sql.DB) {
	statement, err := db.Prepare("SELECT * FROM users")
	if err != nil {
		log.Printf("%v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "select statement went wrong"})
		return
	}
	result, err := statement.Query()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	var users []models.User
	for result.Next() {
		var user models.User
		err := result.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			log.Printf("%v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error scanning user"})
			return
		}
		log.Printf("%v", user)
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}
