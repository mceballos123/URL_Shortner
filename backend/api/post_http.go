package api
import (
	"fmt"
	"time"
	"encoding/json"
	"database/sql"
	"net/http"
)

type CreateUrlRequest struct{
	Alias string `json:"alias" binding:"required"`
	Link string `json:"link" binding:"required"`
	ExpiredAt *time.Time `json:"expired_at"`

}

type CreateUserRequest struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func postCreateUrl(){

}

func postCreateUser(db *sql.DB) http.HandlerFunc{ //Function type for HTTP routes
	return func(w http.ResponseWriter, r *http.Request){ //HTTP Handler function

		var request CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&request)

		if err!=nil{
			fmt.Println("Error decoding request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}else{
			fmt.Println("User created successfully")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))

	}


}