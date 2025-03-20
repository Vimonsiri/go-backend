package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/ployMatsuri/go-backend/config"
	"gitlab.com/ployMatsuri/go-backend/models"
	"gitlab.com/ployMatsuri/go-backend/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config.InitDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := routes.SetupRouter()
	log.Printf("Server is running on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		log.Fatalf("Error querying users table: %v", result.Error)
	} else {
		fmt.Println("Users table exists. Current users count:", result.RowsAffected)
	}
}
