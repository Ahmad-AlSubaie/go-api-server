package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin" //HTTP router
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //posgress driver
)

func setupRouter(db *sql.DB) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	// Creating a route for server.com/ping with responce handeling in the anaonymis funciton
	r.GET("/ping", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.String(503, "pong")
			log.Println("ping, pong DB down:", err)
		} else {
			c.String(http.StatusOK, "pong")
			log.Println("ping, pong DB OK")
		}

	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, err := db.Exec("Select * from users") // no imput sanitization/paramiterization
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
			log.Println("GET request failed:", err)
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	/*
		authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
			"foo":  "bar", // user:foo password:bar
			"manu": "123", // user:manu password:123
		}))
	*/
	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	/*
		authorized.POST("admin", func(c *gin.Context) {
			user := c.MustGet(gin.AuthUserKey).(string)

			// Parse JSON
			var json struct {
				Value string `json:"value" binding:"required"`
			}

			if c.Bind(&json) == nil {
				db[user] = json.Value
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			}
		})
	*/
	return r
}

func setupDatabase() *sql.DB {
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	return db
}
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db := setupDatabase()
	r := setupRouter(db)
	//0.0.0.0:8080
	r.Run(":8080")
}
