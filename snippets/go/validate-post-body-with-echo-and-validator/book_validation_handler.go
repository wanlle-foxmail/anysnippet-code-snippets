package main

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Book struct {
	Name        string `json:"name" validate:"required,min=2,max=80"`
	Page        int    `json:"page" validate:"gte=1,lte=5000"`
	Type        string `json:"type" validate:"required,oneof=fiction nonfiction reference"`
	Description string `json:"description" validate:"required,min=10,max=500"`
}

type CustomValidator struct {
	validator *validator.Validate
}

var ProcessBook = func(book Book) map[string]interface{} {
	return map[string]interface{}{
		"message": "Book request is valid. Continue to the next step.",
		"book":    book,
	}
}

func (cv *CustomValidator) Validate(value interface{}) error {
	return cv.validator.Struct(value)
}

// CreateBookHandler demonstrates a validate-and-continue request flow in Echo.
// Request flow:
//
//	request body -> bind JSON
//	                  |
//	                  +-- invalid JSON ----> return 400
//	                  |
//	                  +-- valid JSON ------> normalize fields
//	                                         |
//	                                         +-- validation fails --> return 400 with friendly field errors
//	                                         |
//	                                         +-- validation passes -> run next step and return 201
func CreateBookHandler(c echo.Context) error {
	var book Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Please send a valid JSON request body.",
		})
	}

	normalizeBook(&book)

	if err := c.Validate(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Please correct the highlighted fields and try again.",
			"errors":  friendlyValidationErrors(err),
		})
	}

	return c.JSON(http.StatusCreated, ProcessBook(book))
}

func NewServer() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/books", CreateBookHandler)
	return e
}

func friendlyValidationErrors(err error) map[string]string {
	messages := map[string]string{}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		messages["request"] = "Please review the submitted fields and try again."
		return messages
	}

	for _, fieldError := range validationErrors {
		messages[strings.ToLower(fieldError.Field())] = friendlyMessageForField(fieldError)
	}

	return messages
}

func friendlyMessageForField(fieldError validator.FieldError) string {
	switch fieldError.Field() {
	case "Name":
		if fieldError.Tag() == "required" {
			return "Please enter a book name."
		}
		return "Keep the name between 2 and 80 characters."
	case "Page":
		if fieldError.Tag() == "lte" {
			return "Please keep the page count at 5000 or fewer."
		}
		return "Please enter a page count greater than 0."
	case "Type":
		return "Choose a valid book type: fiction, nonfiction, or reference."
	case "Description":
		return "Write a description between 10 and 500 characters."
	default:
		return "Please review this field and try again."
	}
}

func normalizeBook(book *Book) {
	book.Name = strings.TrimSpace(book.Name)
	book.Type = strings.TrimSpace(book.Type)
	book.Description = strings.TrimSpace(book.Description)
}

func main() {
	e := NewServer()
	e.Logger.Fatal(e.Start(":8080"))
}
