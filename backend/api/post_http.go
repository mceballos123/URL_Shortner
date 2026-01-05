package api
import (
	"time"
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	return func(c *gin.Context){
		var request CreateUrlRequest
		err := c.ShouldBindJSON(&request) 

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		query := `INSERT INTO urls (alias, url, expires_at, user_id)
		VALUES ($1, $2, $3, $4)
		`
		_, err = db.Exec(query, request.Alias, request.Link, request.ExpiredAt, request.UserId)

		if err !=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(http.StatusCreated, gin.H{"message": "URL created successfully"})
	}

}

func PostCreateUser(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){ 

		var request CreateUserRequest

		err := c.ShouldBind(&request)

		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}
		hashedPassword, err :=bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": "Failed to process password"})
			return 
		}

		query := `INSERT INTO users (username,password)
		VALUES ($1,$2)
		
		`
		_, err = db.Exec(query, request.Username, string(hashedPassword))

		if err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func PostLogin(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		var request LoginRequest
		err := c.BindJSON(&request)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		var storedPassword string
		query :=`SELECT password FROM users WHERE username=$1`
		err = db.QueryRow(query, request.Username).Scan(&storedPassword) 

		if err !=nil{
			if err == sql.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return 
			}else{
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return 
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(request.Password))

		if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}