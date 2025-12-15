package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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


func getAllUrls( db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		query := `SELECT * FROM urls` //Selects the user
		rows, err := db.Query(query)

		if err !=nil{
			fmt.Println("Error querying database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		defer rows.Close() // Schedules the rows to be closed when the function returns

		var urls []Url
		for rows.Next(){
			var u Url 
			err := rows.Scan(&u.Id, &u.Alias, &u.Url, &u.ExpiresAt, &u.UserId)
			if err !=nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 
			}
			urls = append(urls, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urls)
	}
}


func getAllUsers(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		query := `SELECT * FROM users`
		rows, err := db.Query(query)

		if err !=nil {
			fmt.Println("Error querying database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		defer rows.Close() // Schedules the row to close when the function returns

		var users []Username
		for rows.Next(){
			var u Username
			err := rows.Scan(&u.Id, &u.Username, &u.Password)
			if err !=nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 
			}
			users = append(users, u)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		
	}
}

func getUrlsByID(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		idStr := r.URL.Query().Get("id")
		if idStr == ""{
			http.Error(w, "ID is required", http.StatusBadRequest)
			return 
		}
		id, err := strconv.Atoi(idStr)

		if err !=nil{
			fmt.Println("Error converting the id to an integer:", err)
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}
		var u Url

		query := `SELECT * FROM urls WHERE id = $1`

		err = db.QueryRow(query,id).Scan(&u.Id, &u.Alias, &u.Url, &u.ExpiresAt, &u.UserId)

		if err !=nil{
			if err == sql.ErrNoRows{
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}else{
				fmt.Println("Error querying database:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func getUserByID(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		idStr := r.URL.Query().Get("id")
		if idStr == ""{
			http.Error(w, "ID is required", http.StatusBadRequest)
			return 
		}
		id, err := strconv.Atoi(idStr)
		
		if err !=nil{
			fmt.Println("Error converting the id to an integer:", err)
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}
		var u Username
		
		query := `SELECT * FROM users WHERE id = $1`

		err = db.QueryRow(query,id).Scan(&u.Id, &u.Username, &u.Password)

		if err !=nil{
			if err == sql.ErrNoRows{
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}

}