package api
import (
	"fmt"
	"time"
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)
type CreateUrlRequest struct{
	Alias string `json:"alias" binding:"required"`
	Link string `json:"link" binding:"required"`
	ExpiredAt *time.Time `json:"expired_at"`
	UserId int `json:"user_id" binding:"required"`
}

type CreateUserRequest struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func PostCreateUrl( db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){ // r is thge request body
		// w is the response writer, it is used to write the response to the client
		var request CreateUrlRequest
		err := c.ShouldBindJSON(&request) //Maps json to struct
		//This takes the request body and decodes it into the request struct

		if err != nil{
			fmt.Println("Error decoding request body:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}
		fmt.Println("Request body decoded successfully")

		query := `INSERT INTO urls (alias, url, expires_at, user_id)
		VALUES ($1, $2, $3, $4)
		`
		_, err = db.Exec(query, request.Alias, request.Link, request.ExpiredAt, request.UserId)

		if err !=nil{
			fmt.Println("Error executing query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
		fmt.Println("URL created successfully") // 
		c.JSON(http.StatusCreated, gin.H{"message": "URL created successfully"})

	}

}

func PostCreateUser(db *sql.DB) gin.HandlerFunc{ //Function type for HTTP routes
	return func(c *gin.Context){ //HTTP Handler function

		var request CreateUserRequest

		err := c.ShouldBind(&request)

		if err!=nil{
			fmt.Println("Error decoding request body:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}
		fmt.Println("Request body decoded successfully")
		fmt.Println("Username:", request.Username)
		fmt.Println("Password:", request.Password)

		query := `INSERT INTO users (username,password)
		VALUES ($1,$2)
		
		`
		_, err = db.Exec(query, request.Username, request.Password)

		if err !=nil{
			fmt.Println("Error executing query:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}else{
			fmt.Println("User created successfully")
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func PostLogin(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		var request LoginRequest
		err := c.BindJSON(&request)
		// Stores the request body in the request memory adress
		if err != nil{
			fmt.Println("Error decoding request body:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		var storedPassword string
		query :=`SELECT password FROM users WHERE username=$1`
		err = db.QueryRow(query, request.Username).Scan(&storedPassword) 
		// Scaqns the query results into the storedPassword

		if err !=nil{
			if err == sql.ErrNoRows{
				fmt.Println("User not found")
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return 
			}else{
				fmt.Println("Error querying database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return 
			}
		}
		if storedPassword != request.Password{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			fmt.Println("Invalid password")
			return 
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}