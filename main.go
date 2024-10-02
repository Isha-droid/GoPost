package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitializeDatabase sets up the connection to the database using Gorm
func InitializeDatabase() {
	// Load environment variables for DB connection
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Create the connection string for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Connect to the PostgreSQL database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	fmt.Println("Database connection successfully opened")
}

func SetupRoutes(app *fiber.App) {
	// Book CRUD routes
	app.Get("/books", GetAllBooks)       // Read all books
	app.Get("/books/:id", GetBookByID)   // Read a single book by ID
	app.Post("/books", CreateBook)       // Create a new book
	app.Put("/books/:id", UpdateBook)    // Update a book by ID
	app.Delete("/books/:id", DeleteBook) // Delete a book by ID
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	InitializeDatabase()

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	SetupRoutes(app)

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default to port 8081 if not set in .env
	}

	// Start Fiber server
	log.Fatal(app.Listen(":" + port))
}
