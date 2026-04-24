# Echo POST Validation with Go

Validate a `Book` JSON body in Echo with `go-playground/validator` and return friendly field messages instead of raw validator errors.

## What It Does

- Defines a `POST /books` handler in Echo
- Uses a `Book` struct with validation rules for `name`, `page`, `type`, and `description`
- Returns user-facing error messages when validation fails
- Calls the next-step `ProcessBook` function only after validation succeeds

## Book Example

```go
type Book struct {
    Name        string `json:"name" validate:"required,min=2,max=80"`
    Page        int    `json:"page" validate:"gte=1,lte=5000"`
    Type        string `json:"type" validate:"required,oneof=fiction nonfiction reference"`
    Description string `json:"description" validate:"required,min=10,max=500"`
}
```

Allowed `type` values:

- `fiction`
- `nonfiction`
- `reference`

## Usage

```bash
go run book_validation_handler.go
```

Send a valid request:

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Clean Code",
    "page": 464,
    "type": "reference",
    "description": "A practical guide to writing better code."
  }'
```

Success response:

```json
{
  "message": "Book request is valid. Continue to the next step.",
  "book": {
    "name": "Clean Code",
    "page": 464,
    "type": "reference",
    "description": "A practical guide to writing better code."
  }
}
```

Validation failure response:

```json
{
  "message": "Please correct the highlighted fields and try again.",
  "errors": {
    "name": "Please enter a book name.",
    "page": "Please enter a page count greater than 0.",
    "type": "Choose a valid book type: fiction, nonfiction, or reference.",
    "description": "Write a description between 10 and 500 characters."
  }
}
```

## Notes

- The handler trims surrounding whitespace from `name`, `type`, and `description` before validation.
- `ProcessBook` is the placeholder for the next step after validation passes.
- Malformed JSON gets a generic user-friendly message instead of a raw bind error.

## Verification

```bash
go mod tidy
go test -race ./...
```