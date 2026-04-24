package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")
var ErrProductNotFound = errors.New("product not found")

const productCacheTTL = 5 * time.Minute

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type productResponse struct {
	Message string  `json:"message"`
	Source  string  `json:"source,omitempty"`
	Product Product `json:"product,omitempty"`
}

type errorResponse struct {
	Message string `json:"message"`
}

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var sampleProducts = map[string]Product{
	"p-100": {
		ID:       "p-100",
		Name:     "Mechanical Keyboard",
		Price:    129.99,
		Category: "hardware",
	},
	"p-200": {
		ID:       "p-200",
		Name:     "USB-C Dock",
		Price:    89.50,
		Category: "accessories",
	},
}

var CacheGet = func(ctx context.Context, key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCacheMiss
	}
	if err != nil {
		return "", fmt.Errorf("read redis cache: %w", err)
	}
	return value, nil
}

var CacheSet = func(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := redisClient.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("write redis cache: %w", err)
	}
	return nil
}

var FindProductByID = func(_ context.Context, id string) (Product, error) {
	product, ok := sampleProducts[id]
	if !ok {
		return Product{}, ErrProductNotFound
	}
	return product, nil
}

func productCacheKey(id string) string {
	return "product:" + id
}

func ReadProductFromCache(ctx context.Context, id string) (Product, error) {
	rawValue, err := CacheGet(ctx, productCacheKey(id))
	if err != nil {
		return Product{}, err
	}

	var product Product
	if err := json.Unmarshal([]byte(rawValue), &product); err != nil {
		return Product{}, fmt.Errorf("decode cached product: %w", err)
	}

	return product, nil
}

func WriteProductToCache(ctx context.Context, product Product) error {
	rawValue, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("encode product for cache: %w", err)
	}

	if err := CacheSet(ctx, productCacheKey(product.ID), string(rawValue), productCacheTTL); err != nil {
		return err
	}

	return nil
}

// GetProductHandler demonstrates a cache-aside lookup in Echo.
// Cache-aside flow:
//
//	request -> try Redis
//	             |
//	             +-- hit ----------> return cached product
//	             |
//	             +-- miss/error ---> load product from backing store
//	                                  |
//	                                  +-- not found ---> return 404
//	                                  |
//	                                  +-- found -------> write Redis best-effort
//	                                                     |
//	                                                     +-> return fresh product
func GetProductHandler(c echo.Context) error {
	productID := strings.TrimSpace(c.Param("id"))
	if productID == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "Please provide a product id."})
	}

	ctx := c.Request().Context()
	cachedProduct, err := ReadProductFromCache(ctx, productID)
	if err == nil {
		return c.JSON(http.StatusOK, productResponse{
			Message: "Product loaded from cache.",
			Source:  "cache",
			Product: cachedProduct,
		})
	}

	if !errors.Is(err, ErrCacheMiss) {
		c.Logger().Warnf("cache read failed for %s: %v", productID, err)
	}

	product, err := FindProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{Message: "Product not found."})
		}

		c.Logger().Errorf("database lookup failed for %s: %v", productID, err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "Unable to load the product right now."})
	}

	if err := WriteProductToCache(ctx, product); err != nil {
		c.Logger().Warnf("cache write failed for %s: %v", productID, err)
	}

	return c.JSON(http.StatusOK, productResponse{
		Message: "Product loaded from database and cached.",
		Source:  "database",
		Product: product,
	})
}

func NewServer() *echo.Echo {
	e := echo.New()
	e.GET("/products/:id", GetProductHandler)
	return e
}

func main() {
	server := NewServer()
	server.Logger.Fatal(server.Start(":8080"))
}
