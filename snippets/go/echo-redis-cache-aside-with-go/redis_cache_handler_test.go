package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func sampleProduct() Product {
	return Product{
		ID:       "p-100",
		Name:     "Mechanical Keyboard",
		Price:    129.99,
		Category: "hardware",
	}
}

func getProductResponse(t *testing.T, id string) (int, map[string]interface{}) {
	t.Helper()

	e := NewServer()
	req := httptest.NewRequest(http.MethodGet, "/products/"+url.PathEscape(id), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Body.Len() == 0 {
		return rec.Code, nil
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	return rec.Code, body
}

func assertProductBody(t *testing.T, body map[string]interface{}, product Product) {
	t.Helper()

	item, ok := body["product"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing product payload: %#v", body)
	}

	if item["id"] != product.ID {
		t.Fatalf("expected id %q, got %#v", product.ID, item["id"])
	}
	if item["name"] != product.Name {
		t.Fatalf("expected name %q, got %#v", product.Name, item["name"])
	}
	if item["category"] != product.Category {
		t.Fatalf("expected category %q, got %#v", product.Category, item["category"])
	}
	if item["price"] != product.Price {
		t.Fatalf("expected price %v, got %#v", product.Price, item["price"])
	}
}

func TestGetProduct_CacheHit(t *testing.T) {
	product := sampleProduct()
	var dbCalls int
	var cacheSetCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, key string) (string, error) {
		if key != productCacheKey(product.ID) {
			t.Fatalf("unexpected cache key: %s", key)
		}

		payload, err := json.Marshal(product)
		if err != nil {
			t.Fatalf("marshal cached product: %v", err)
		}

		return string(payload), nil
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		cacheSetCalls++
		return nil
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		dbCalls++
		return Product{}, nil
	}

	statusCode, body := getProductResponse(t, product.ID)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
	if body["source"] != "cache" {
		t.Fatalf("expected cache source, got %#v", body["source"])
	}
	if body["message"] != "Product loaded from cache." {
		t.Fatalf("unexpected message: %#v", body["message"])
	}
	if dbCalls != 0 {
		t.Fatalf("expected database not to be called, got %d calls", dbCalls)
	}
	if cacheSetCalls != 0 {
		t.Fatalf("expected cache set not to be called, got %d calls", cacheSetCalls)
	}

	assertProductBody(t, body, product)
}

func TestGetProduct_CacheMissLoadsDatabaseAndBackfillsCache(t *testing.T) {
	product := sampleProduct()
	var dbCalls int
	var cacheSetCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, key string) (string, error) {
		if key != productCacheKey(product.ID) {
			t.Fatalf("unexpected cache key: %s", key)
		}
		return "", ErrCacheMiss
	}
	CacheSet = func(_ context.Context, key string, value string, ttl time.Duration) error {
		cacheSetCalls++
		if key != productCacheKey(product.ID) {
			t.Fatalf("unexpected cache key: %s", key)
		}
		if ttl != productCacheTTL {
			t.Fatalf("unexpected ttl: %v", ttl)
		}

		var cached Product
		if err := json.Unmarshal([]byte(value), &cached); err != nil {
			t.Fatalf("invalid cached product JSON: %v", err)
		}
		if cached.ID != product.ID {
			t.Fatalf("unexpected cached product id: %s", cached.ID)
		}

		return nil
	}
	FindProductByID = func(_ context.Context, id string) (Product, error) {
		dbCalls++
		if id != product.ID {
			t.Fatalf("unexpected database id: %s", id)
		}
		return product, nil
	}

	statusCode, body := getProductResponse(t, product.ID)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
	if body["source"] != "database" {
		t.Fatalf("expected database source, got %#v", body["source"])
	}
	if body["message"] != "Product loaded from database and cached." {
		t.Fatalf("unexpected message: %#v", body["message"])
	}
	if dbCalls != 1 {
		t.Fatalf("expected one database call, got %d", dbCalls)
	}
	if cacheSetCalls != 1 {
		t.Fatalf("expected one cache write, got %d", cacheSetCalls)
	}

	assertProductBody(t, body, product)
}

func TestGetProduct_CacheMissReturnsNotFound(t *testing.T) {
	var cacheSetCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, _ string) (string, error) {
		return "", ErrCacheMiss
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		cacheSetCalls++
		return nil
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		return Product{}, ErrProductNotFound
	}

	statusCode, body := getProductResponse(t, "missing-product")

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
	if body["message"] != "Product not found." {
		t.Fatalf("unexpected message: %#v", body["message"])
	}
	if cacheSetCalls != 0 {
		t.Fatalf("expected cache write to be skipped, got %d calls", cacheSetCalls)
	}
}

func TestGetProduct_CacheMissReturnsServerErrorWhenDatabaseFails(t *testing.T) {
	var cacheSetCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, _ string) (string, error) {
		return "", ErrCacheMiss
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		cacheSetCalls++
		return nil
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		return Product{}, errors.New("database offline")
	}

	statusCode, body := getProductResponse(t, "p-100")

	if statusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", statusCode)
	}
	if body["message"] != "Unable to load the product right now." {
		t.Fatalf("unexpected message: %#v", body["message"])
	}
	if cacheSetCalls != 0 {
		t.Fatalf("expected cache write to be skipped, got %d calls", cacheSetCalls)
	}
}

func TestGetProduct_CacheErrorFallsBackToDatabase(t *testing.T) {
	product := sampleProduct()
	var dbCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, _ string) (string, error) {
		return "", errors.New("redis unavailable")
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		return nil
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		dbCalls++
		return product, nil
	}

	statusCode, body := getProductResponse(t, product.ID)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
	if body["source"] != "database" {
		t.Fatalf("expected database source, got %#v", body["source"])
	}
	if dbCalls != 1 {
		t.Fatalf("expected one database call, got %d", dbCalls)
	}

	assertProductBody(t, body, product)
}

func TestGetProduct_CacheWriteFailureDoesNotBreakResponse(t *testing.T) {
	product := sampleProduct()
	var dbCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, _ string) (string, error) {
		return "", ErrCacheMiss
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		return errors.New("cache write failed")
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		dbCalls++
		return product, nil
	}

	statusCode, body := getProductResponse(t, product.ID)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
	if body["source"] != "database" {
		t.Fatalf("expected database source, got %#v", body["source"])
	}
	if dbCalls != 1 {
		t.Fatalf("expected one database call, got %d", dbCalls)
	}

	assertProductBody(t, body, product)
}

func TestGetProduct_InvalidProductIDReturnsBadRequest(t *testing.T) {
	var cacheGetCalls int
	var dbCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, _ string) (string, error) {
		cacheGetCalls++
		return "", ErrCacheMiss
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		return nil
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		dbCalls++
		return Product{}, nil
	}

	statusCode, body := getProductResponse(t, "   ")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
	if body["message"] != "Please provide a product id." {
		t.Fatalf("unexpected message: %#v", body["message"])
	}
	if cacheGetCalls != 0 {
		t.Fatalf("expected cache read to be skipped, got %d calls", cacheGetCalls)
	}
	if dbCalls != 0 {
		t.Fatalf("expected database read to be skipped, got %d calls", dbCalls)
	}
}

func TestGetProduct_InvalidCachedJSONFallsBackToDatabase(t *testing.T) {
	product := sampleProduct()
	var dbCalls int
	var cacheSetCalls int

	originalCacheGet := CacheGet
	originalCacheSet := CacheSet
	originalFindProductByID := FindProductByID
	defer func() {
		CacheGet = originalCacheGet
		CacheSet = originalCacheSet
		FindProductByID = originalFindProductByID
	}()

	CacheGet = func(_ context.Context, _ string) (string, error) {
		return "{not-valid-json}", nil
	}
	CacheSet = func(_ context.Context, _ string, _ string, _ time.Duration) error {
		cacheSetCalls++
		return nil
	}
	FindProductByID = func(_ context.Context, _ string) (Product, error) {
		dbCalls++
		return product, nil
	}

	statusCode, body := getProductResponse(t, product.ID)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
	if body["source"] != "database" {
		t.Fatalf("expected database source, got %#v", body["source"])
	}
	if dbCalls != 1 {
		t.Fatalf("expected one database call, got %d", dbCalls)
	}
	if cacheSetCalls != 1 {
		t.Fatalf("expected one cache write, got %d", cacheSetCalls)
	}

	assertProductBody(t, body, product)
}
