package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB{
	err := godotenv.Load()

	if err != nil{
		log.Fatal("Error loading .env file")
	}

	postGresUsername := os.Getenv("POSTGRES_USER")
	postGresPassword := os.Getenv("POSTGRES_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, postGresUsername, postGresPassword, DBName)

	db, err := sql.Open("postgres",connStr)
	if err !=nil{
		log.Fatal("Failed to connect to the database:", err)
	}else{
		log.Println("Connected to the database successfully")
	}

	defer db.Close()

	err = db.Ping()
	
	if err != nil{
		log.Fatal("Failed to ping the database:", err)
	}else{
		log.Println("Pinged the database successfully")
	}
	return db 
}