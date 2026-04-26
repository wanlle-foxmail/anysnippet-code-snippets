package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Pagination struct {
	Page     int
	PageSize int
	Offset   int
}

// Flow: read page and page_size -> apply defaults and validation -> cap the page size -> return page, size, and offset.
func ParsePaginationQueryParams(values url.Values, defaultPage, defaultPageSize, maxPageSize int) (Pagination, error) {
	if defaultPage <= 0 {
		return Pagination{}, errors.New("default page must be greater than 0")
	}
	if defaultPageSize <= 0 {
		return Pagination{}, errors.New("default page size must be greater than 0")
	}
	if maxPageSize <= 0 {
		return Pagination{}, errors.New("max page size must be greater than 0")
	}

	page, err := parsePositiveIntOrDefault(values.Get("page"), defaultPage)
	if err != nil {
		return Pagination{}, fmt.Errorf("parse page: %w", err)
	}

	pageSize, err := parsePositiveIntOrDefault(values.Get("page_size"), defaultPageSize)
	if err != nil {
		return Pagination{}, fmt.Errorf("parse page size: %w", err)
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}, nil
}

func parsePositiveIntOrDefault(raw string, fallback int) (int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, errors.New("must be an integer")
	}
	if value <= 0 {
		return 0, errors.New("must be greater than 0")
	}

	return value, nil
}

func main() {
	values := url.Values{
		"page":      {"2"},
		"page_size": {"100"},
	}

	pagination, err := ParsePaginationQueryParams(values, 1, 25, 50)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("page=%d page_size=%d offset=%d", pagination.Page, pagination.PageSize, pagination.Offset)
}
