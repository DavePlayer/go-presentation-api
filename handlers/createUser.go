package handlers

import (
	"apitest/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context, db *sql.DB) {
	var newUser models.NewUser
	if err := c.BindJSON(&newUser); err != nil {
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

	// INSERT operation
	statement, err := tx.Prepare("INSERT INTO users (login, password) VALUES (?, ?)")
	if err != nil {
		log.Printf("insert prepare error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "insert prepare error"})
		return
	}
	_, err = statement.Exec(newUser.Login, newUser.Password)
	if err != nil {
		log.Printf("insert execute error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not insert user"})
		return
	}

	// SELECT operation
	statement, err = tx.Prepare("SELECT id, login, password FROM users WHERE login = ? AND password = ?")
	if err != nil {
		log.Printf("select prepare error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "select prepare error"})
		return
	}

	result, err := statement.Query(newUser.Login, newUser.Password)
	if err != nil {
		log.Printf("select execute error: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "can't return user when fetched back"})
		return
	}

	if result.Next() {
		var user models.User
		err := result.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			log.Printf("select scan error: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error scanning user"})
			return
		}
		log.Printf("%v", user)
		c.IndentedJSON(http.StatusOK, user)

		// Commit the transaction after successful insert and select
		err = tx.Commit()
		if err != nil {
			log.Printf("transaction commit error: %v", err)
			// Handle commit error if needed
		}
		return
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no new user returned back"})
		return
	}

}
