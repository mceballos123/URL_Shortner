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


func postCreateUrl( db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){ // r is thge request body
		// w is the response writer, it is used to write the response to the client
		var request CreateUrlRequest
		err := json.NewDecoder(r.Body).Decode(&request) //Maps json to struct
		//This takes the request body and decodes it into the request struct

		if err != nil{
			fmt.Println("Error decoding request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		fmt.Println("Request body decoded successfully")

		query := `INSERT INTO urls (alias, url, expires_at, user_id)
		VALUES ($1, $2, $3, $4)
		`
		_, err = db.Exec(query, request.Alias, request.Link, request.ExpiredAt)

		if err !=nil{
			fmt.Println("Error executing query:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		fmt.Println("URL created successfully") // 
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("URL created"))

	}

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

func postLogin(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var request LoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil{
			fmt.Println("Error decoding request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		var storedPassword string
		query :=`SELECT password FROM users WHERE username=$1`
		err = db.QueryRow(query, request.Username).Scan(&storedPassword)

		if err !=nil{
			if err == sql.ErrNoRows{
				fmt.Println("User not found")
				http.Error(w, "User not found", http.StatusNotFound)
				return 
			}else{
				fmt.Println("Error querying database:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 
			}

			
		}
		if storedPassword != request.Password{
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			fmt.Println("Invalid password")
			return 
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	}
}