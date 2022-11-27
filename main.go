package main

import (
	"log"
	"net/http"
	"os"
	"postgres-gorm/models"
	"postgres-gorm/storage"

	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBooks(c *fiber.Ctx) {
	var book Book // can also be: book := Book{}

	// * TIP: Parse the request body into the book struct as go does understand json
	err := c.BodyParser(&book)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"error": err.Error()})
		return
		// return err
	}

	err1 := r.DB.Create(&book)

	if err1 != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "Book cannot be created"})
		return
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book created successfully"})

	// return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) {
	bookModels := &[]models.Book{}

	err := r.DB.Find(&bookModels)

	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "Books cannot be fetched"})
		return
	}

	// fmt.Println(len(*bookModels))

	if len(*bookModels) == 0 {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "No books found"})
		return
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books fetched successfully", "books": bookModels})

	// return nil
}

func (r *Repository) DeleteBooks(c *fiber.Ctx) {
	bookID := c.Params("id")

	bookModel := &models.Book{}

	err := r.DB.Find(&bookModel, bookID).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "Book cannot be found"})
		return
	}

	err1 := r.DB.Delete(&bookModel, bookID).Error

	if err1 != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "Book cannot be deleted"})
		return
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book deleted successfully"})

	// return nil
}

func (r *Repository) GetBookByID(c *fiber.Ctx) {
	bookID := c.Params("id")

	bookModel := &models.Book{}

	err := r.DB.Where("id = ?", bookID).First(&bookModel)

	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "Book cannot be found"})
		return
	}

	if bookModel.ID == 0 {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "No book found"})
		return
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book fetched successfully", "book": bookModel})

	// return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/create_books", r.CreateBooks)
	api.Delete("/delete_books/:id", r.DeleteBooks)
	api.Get("/get_book/:id", r.GetBookByID)
	api.Get("/get_books", r.GetBooks)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err1 := models.MigrateBooks(db)

	if err1 != nil {
		log.Fatal("Error migrating database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(3000)
}
