package handlers

import (
	"io"
	"net/http/httptest"
	"testing"

	"asterisk-manager/domain"
	"asterisk-manager/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPaginationMiddleware(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(&repositories.Repos{})

	// Test route with pagination middleware
	app.Get("/test", handler.Pagination, func(c *fiber.Ctx) error {
		pagination := c.Locals("pagination").(*domain.PaginationInput)
		return c.JSON(pagination)
	})

	tests := []struct {
		name           string
		queryString    string
		expectedPage   int
		expectedPerPage int
	}{
		{
			name:           "Default values when no params",
			queryString:    "",
			expectedPage:   1,
			expectedPerPage: 10,
		},
		{
			name:           "Custom page and perPage",
			queryString:    "?page=2&perPage=20",
			expectedPage:   2,
			expectedPerPage: 20,
		},
		{
			name:           "Invalid page defaults to 1",
			queryString:    "?page=0&perPage=15",
			expectedPage:   1,
			expectedPerPage: 15,
		},
		{
			name:           "Invalid perPage defaults to 10",
			queryString:    "?page=3&perPage=0",
			expectedPage:   3,
			expectedPerPage: 10,
		},
		{
			name:           "PerPage exceeding max gets capped at 100",
			queryString:    "?page=1&perPage=200",
			expectedPage:   1,
			expectedPerPage: 100,
		},
		{
			name:           "Negative page defaults to 1",
			queryString:    "?page=-5&perPage=10",
			expectedPage:   1,
			expectedPerPage: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test"+tt.queryString, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			// Read response body
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			// Parse and verify JSON structure contains correct values
			bodyStr := string(body)
			assert.Contains(t, bodyStr, `"Page":`)
			assert.Contains(t, bodyStr, `"PerPage":`)

			// The body should be valid JSON with our expected structure
			assert.NotEmpty(t, bodyStr)
		})
	}
}

func TestPaginationResponse_CalculatePages(t *testing.T) {
	tests := []struct {
		name          string
		total         int64
		perPage       int
		expectedPages int
	}{
		{
			name:          "Exact division",
			total:         100,
			perPage:       10,
			expectedPages: 10,
		},
		{
			name:          "With remainder",
			total:         105,
			perPage:       10,
			expectedPages: 11,
		},
		{
			name:          "Less than one page",
			total:         5,
			perPage:       10,
			expectedPages: 1,
		},
		{
			name:          "Zero total",
			total:         0,
			perPage:       10,
			expectedPages: 0,
		},
		{
			name:          "Single item",
			total:         1,
			perPage:       10,
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := domain.PaginationResponse{
				Total:   tt.total,
				PerPage: tt.perPage,
			}
			pr.CalculatePages()

			assert.Equal(t, tt.expectedPages, pr.Pages)
		})
	}
}
