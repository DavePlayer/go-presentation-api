package handlers

import (
	"apitest/models"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserByID(c *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("%v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "select statement went wrong"})
		return
	}
	statement, err := db.Prepare("SELECT id, login, password FROM users WHERE id = ?")
	if err != nil {
		log.Printf("%v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "select statement went wrong"})
		return
	}
	result, err := statement.Query(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if result.Next() {
		var user models.User
		err := result.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			log.Printf("%v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error scanning user"})
			return
		}
		log.Printf("%v", user)
		c.IndentedJSON(http.StatusOK, user)
		return
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
}
