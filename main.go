package main

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"log"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserToUpdate struct {
	Id   int  `json:"id"`
	User User `json:"User"`
}

func test(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func UpdateUser(c *gin.Context, db *sql.DB) {
	var userData UserToUpdate
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
	var user User
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

func userByID(c *gin.Context, db *sql.DB) {
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
		var user User
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

func getUsers(c *gin.Context, db *sql.DB) {
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
	var users []User
	for result.Next() {
		var user User
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

func createUser(c *gin.Context, db *sql.DB) {
	var newUser NewUser
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
		var user User
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

func initDb(db *sql.DB) error {
	statement, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT,
		password TEXT
	);
	`)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	users := []NewUser{{"user1", "secret"}, {"user2", "secret"}}
	for _, user := range users {
		log.Printf("inserting %v\n", user)
		statement, err := db.Prepare("INSERT INTO users (login, password) VALUES (?, ?)")
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = statement.Exec(user.Login, user.Password)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

func main() {
	gin.SetMode(gin.DebugMode)
	database, dbErr := sql.Open("sqlite3", "./test.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	initDb(database)
	router := gin.Default()

	router.GET("/", test)
	router.GET("/users/", func(c *gin.Context) {
		getUsers(c, database)
	})
	router.GET("/users/:id", func(c *gin.Context) {
		userByID(c, database)
	})
	router.POST("/users/add", func(c *gin.Context) {
		createUser(c, database)
	})
	router.PATCH("/users/update", func(c *gin.Context) {
		UpdateUser(c, database)
	})

	router.Run("127.0.0.1:8080")
}
