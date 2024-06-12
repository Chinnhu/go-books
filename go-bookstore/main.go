package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go-bookstore/models"
	"go-bookstore/storage"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Year   string `json:"year"`
	ISBN   string `json:"isbn"`
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "error: request failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "error creating book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been added to the database"})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Book{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "couldn't get books"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books received successfully", "data": bookModels})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := &models.Book{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cant be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "couldnt delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been deleted from the database"})
	return nil

}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Book{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id field cant be empty"})
		return nil
	}
	fmt.Println("The ID is ", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "couldnt get book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book id found successfully", "data": bookModel})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/createBook", r.CreateBook)
	api.Delete("deleteBook/:id", r.DeleteBook)
	api.Get("/getBook/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("failed to load the db")
	}
	err = models.MigrateBook(db)
	if err != nil {
		log.Fatal("couldnt migrate db")
	}

	r := Repository{DB: db}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
