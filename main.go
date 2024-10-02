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

// Book model
type Book struct {
	gorm.Model
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

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

	// Migrate the schema (create table if it doesn't exist)
	DB.AutoMigrate(&Book{}) // Migrate the Book model

	fmt.Println("Database connection successfully opened and Book table migrated")
}

// GetAllBooks retrieves all books from the database
func GetAllBooks(c *fiber.Ctx) error {
	var books []Book
	if result := DB.Find(&books); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve books",
		})
	}
	return c.JSON(books)
}

// GetBookByID retrieves a specific book by ID
func GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if result := DB.First(&book, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Book not found",
		})
	}
	return c.JSON(book)
}

// CreateBook creates a new book
func CreateBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	if result := DB.Create(book); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}
	return c.JSON(book)
}

// UpdateBook updates an existing book by ID
func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if result := DB.First(&book, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if result := DB.Save(&book); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}
	return c.JSON(book)
}

// DeleteBook deletes a book by ID
func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if result := DB.First(&book, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	if result := DB.Delete(&book); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// SetupRoutes configures the routes for the Fiber app
func SetupRoutes(app *fiber.App) {
	app.Get("/books", GetAllBooks)
	app.Get("/books/:id", GetBookByID)
	app.Post("/books", CreateBook)
	app.Put("/books/:id", UpdateBook)
	app.Delete("/books/:id", DeleteBook)
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
