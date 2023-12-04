package handlers

import (
	"apitest/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context, db *sql.DB) {
	var userData models.UserToUpdate
	if err := c.BindJSON(&userData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data passed"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("transaction begin error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "transaction begin error"})
		return
	}
	defer tx.Rollback()

	// UPDATE operation
	_, err = tx.Exec("UPDATE users SET login = ?, password = ? WHERE id = ?", userData.User.Login, userData.User.Password, userData.Id)
	if err != nil {
		log.Printf("update execute error: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "could not update user"})
		return
	}

	// SELECT operation
	var user models.User
	err = tx.QueryRow("SELECT id, login, password FROM users WHERE id = ?", userData.Id).Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		log.Printf("select scan error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error scanning user"})
		return
	}

	log.Printf("%v", user)
	c.IndentedJSON(http.StatusOK, user)

	err = tx.Commit()
	if err != nil {
		log.Printf("transaction commit error: %v", err)
		// Handle commit error if needed
	}
}
