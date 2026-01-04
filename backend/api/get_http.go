package api
import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)
type Url struct{
	Id int `json:"id"`
	Alias string `json:"alias"`
	Url string `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	UserId int `json:"user_id"`
}
type Username struct{
	Id int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetAllUrls( db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		query := `SELECT * FROM urls` //Selects the user
		rows, err := db.Query(query)

		if err !=nil{
			fmt.Println("Error querying database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
		defer rows.Close() // Schedules the rows to be closed when the function returns

		var urls []Url
		for rows.Next(){
			var u Url 
			err := rows.Scan(&u.Id, &u.Alias, &u.Url, &u.ExpiresAt, &u.UserId)
			if err !=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return 
			}
			urls = append(urls, u)
		}
		c.JSON(http.StatusOK, urls)
		
	}
}

func GetAllUsers(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		query := `SELECT * FROM users`
		rows, err := db.Query(query)

		if err !=nil {
			fmt.Println("Error querying database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
		defer rows.Close() // Schedules the row to close when the function returns

		var users []Username
		for rows.Next(){
			var u Username
			err := rows.Scan(&u.Id, &u.Username, &u.Password)
			if err !=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return 
			}
			users = append(users, u)
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUrlsByID(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		idStr := c.Query("id")
		if idStr == ""{
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return 
		}
		id, err := strconv.Atoi(idStr)

		if err !=nil{
			fmt.Println("Error converting the id to an integer:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
			return
		}
		var u Url

		query := `SELECT * FROM urls WHERE id = $1`

		err = db.QueryRow(query,id).Scan(&u.Id, &u.Alias, &u.Url, &u.ExpiresAt, &u.UserId)

		if err !=nil{
			if err == sql.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			}else{
				fmt.Println("Error querying database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, u)
	}
}

func GetUserByID(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		idStr := c.Query("id")
		if idStr == ""{
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return 
		}
		id, err := strconv.Atoi(idStr)
		
		if err !=nil{
			fmt.Println("Error converting the id to an integer:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
			return
		}
		var u Username
		
		query := `SELECT * FROM users WHERE id = $1`

		err = db.QueryRow(query,id).Scan(&u.Id, &u.Username, &u.Password)

		if err !=nil{
			if err == sql.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
		}
		c.JSON(http.StatusOK, u)
	}
}

func RedirectShortUrl(db *sql.DB) gin.HandlerFunc{
	return func(c  *gin.Context){
		alias := c.Param("alias") // When the user visits /google it will take them to google.com

		if alias == "" {
			fmt.Println("Alias is either incorrect or user didn't insert a alias")
			c.JSON(http.StatusBadRequest, gin.H{"error":"Alias is required"})
			return
		}

		query := `SELECT url, expires_at FROM WHERE alias =$1`

		var fullUrl string
		var expiresAt sql.NullTime
		err := db.QueryRow(query,alias).Scan(&fullUrl, &expiresAt)

		if err != nil{
			if err == sql.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error":err})
				return 
			}else{
				fmt.Println("Error with quering the db: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return 
			}
		}
		
		if expiresAt.Valid && time.Now().After(expiresAt.Time){
			c.JSON(http.StatusGone, gin.H{"error":"This short URL has expired"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect,fullUrl)
		
	}
}