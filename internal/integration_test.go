package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ooliokartchallenge/internal/application/services"
	"ooliokartchallenge/internal/domain/entities"
	httpInfra "ooliokartchallenge/internal/infrastruture/http"
	"ooliokartchallenge/internal/infrastruture/http/handlers"
	"ooliokartchallenge/internal/infrastruture/http/middleware"
	"ooliokartchallenge/internal/infrastruture/repositories"
	"ooliokartchallenge/pkg/logger"
)

// TestServer holds the test server and dependencies
type TestServer struct {
	server  *httptest.Server
	handler http.Handler
}

// setupTestServer creates a test server with all dependencies
func setupTestServer(t *testing.T) *TestServer {
	// Initialize logger
	appLogger := logger.New()

	// Initialize repositories
	productRepo := repositories.NewProductRepository()

	// Create test coupon files for promo repository
	couponFiles := []string{
		"../testdata/couponbase1.txt",
		"../testdata/couponbase2.txt",
		"../testdata/couponbase3.txt",
	}
	promoRepo := repositories.NewPromoRepository(couponFiles)

	// Initialize services
	promoService := services.NewPromoService(promoRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(productRepo, promoService)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, appLogger)
	orderHandler := handlers.NewOrderHandler(orderService, appLogger)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(appLogger)
	corsMiddleware := middleware.NewCORSMiddleware()

	// Initialize router
	router := httpInfra.NewRouter(productHandler, orderHandler, authMiddleware, corsMiddleware)
	handler := router.SetupRoutes()

	// Create test server
	server := httptest.NewServer(handler)

	return &TestServer{
		server:  server,
		handler: handler,
	}
}

// TestOpenAPICompliance tests all endpoints for OpenAPI specification compliance
func TestOpenAPICompliance(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.server.Close()

	t.Run("Product Endpoints", func(t *testing.T) {
		testProductEndpoints(t, testServer)
	})

	t.Run("Order Endpoints", func(t *testing.T) {
		testOrderEndpoints(t, testServer)
	})

	t.Run("Error Response Format", func(t *testing.T) {
		testErrorResponseFormat(t, testServer)
	})

	t.Run("Content-Type Headers", func(t *testing.T) {
		testContentTypeHeaders(t, testServer)
	})
}

// testProductEndpoints validates product endpoint compliance
func testProductEndpoints(t *testing.T, testServer *TestServer) {
	t.Run("GET /product - List all products", func(t *testing.T) {
		resp, err := http.Get(testServer.server.URL + "/product")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Verify status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verify Content-Type header
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		// Verify response schema
		var products []entities.Product
		if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Validate product schema compliance
		for i, product := range products {
			validateProductSchema(t, product, fmt.Sprintf("product[%d]", i))
		}
	})

	t.Run("GET /product/{productId} - Get specific product", func(t *testing.T) {
		// Test with valid product ID (using ID "10" which exists in sample data)
		resp, err := http.Get(testServer.server.URL + "/product/10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Verify status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verify Content-Type header
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		// Verify response schema
		var product entities.Product
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		validateProductSchema(t, product, "product")
	})

	t.Run("GET /product/{productId} - Invalid product ID", func(t *testing.T) {
		resp, err := http.Get(testServer.server.URL + "/product/invalid")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		validateErrorResponse(t, resp)
	})

	t.Run("GET /product/{productId} - Non-existent product", func(t *testing.T) {
		resp, err := http.Get(testServer.server.URL + "/product/999")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		validateErrorResponse(t, resp)
	})
}

func testOrderEndpoints(t *testing.T, testServer *TestServer) {
	t.Run("POST /order - Valid order without auth", func(t *testing.T) {
		orderRequest := entities.OrderRequest{
			Items: []entities.OrderItem{
				{ProductID: "10", Quantity: 2},
			},
		}

		body, _ := json.Marshal(orderRequest)
		resp, err := http.Post(testServer.server.URL+"/order", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}

		validateErrorResponse(t, resp)
	})

	t.Run("POST /order - Valid order with auth", func(t *testing.T) {
		orderRequest := entities.OrderRequest{
			Items: []entities.OrderItem{
				{ProductID: "10", Quantity: 2},
			},
		}

		body, _ := json.Marshal(orderRequest)
		req, _ := http.NewRequest("POST", testServer.server.URL+"/order", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api_key", "apitest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		var order entities.Order
		if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		validateOrderSchema(t, order)
	})

	t.Run("POST /order - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", testServer.server.URL+"/order", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api_key", "apitest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		validateErrorResponse(t, resp)
	})

	t.Run("POST /order - Empty items", func(t *testing.T) {
		orderRequest := entities.OrderRequest{
			Items: []entities.OrderItem{},
		}

		body, _ := json.Marshal(orderRequest)
		req, _ := http.NewRequest("POST", testServer.server.URL+"/order", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api_key", "apitest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		validateErrorResponse(t, resp)
	})
}

func testErrorResponseFormat(t *testing.T, testServer *TestServer) {
	t.Run("404 Not Found format", func(t *testing.T) {
		resp, err := http.Get(testServer.server.URL + "/nonexistent")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

	})
}

// testContentTypeHeaders validates Content-Type headers for all responses
func testContentTypeHeaders(t *testing.T, testServer *TestServer) {
	testCases := []struct {
		name           string
		method         string
		path           string
		body           string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:           "GET /product",
			method:         "GET",
			path:           "/product",
			expectedStatus: 200,
		},
		{
			name:           "GET /product/1",
			method:         "GET",
			path:           "/product/1",
			expectedStatus: 200,
		},
		{
			name:           "POST /order with auth",
			method:         "POST",
			path:           "/order",
			body:           `{"items":[{"productId":"10","quantity":1}]}`,
			headers:        map[string]string{"api_key": "apitest"},
			expectedStatus: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.body != "" {
				req, err = http.NewRequest(tc.method, testServer.server.URL+tc.path, strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tc.method, testServer.server.URL+tc.path, nil)
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}


			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()


			contentType := resp.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}
		})
	}
}

func validateProductSchema(t *testing.T, product entities.Product, context string) {
	if product.ID == "" {
		t.Errorf("%s: missing required field 'id'", context)
	}
	if product.Name == "" {
		t.Errorf("%s: missing required field 'name'", context)
	}
	if product.Price <= 0 {
		t.Errorf("%s: invalid price value: %f", context, product.Price)
	}
	if product.Category == "" {
		t.Errorf("%s: missing required field 'category'", context)
	}

	if product.Image.Thumbnail == "" {
		t.Errorf("%s: missing required field 'image.thumbnail'", context)
	}
	if product.Image.Mobile == "" {
		t.Errorf("%s: missing required field 'image.mobile'", context)
	}
	if product.Image.Tablet == "" {
		t.Errorf("%s: missing required field 'image.tablet'", context)
	}
	if product.Image.Desktop == "" {
		t.Errorf("%s: missing required field 'image.desktop'", context)
	}
}

func validateOrderSchema(t *testing.T, order entities.Order) {
	if order.ID == "" {
		t.Error("Order: missing required field 'id'")
	}
	if order.Total < 0 {
		t.Errorf("Order: invalid total value: %f", order.Total)
	}
	if order.Discounts < 0 {
		t.Errorf("Order: invalid discounts value: %f", order.Discounts)
	}
	if len(order.Items) == 0 {
		t.Error("Order: missing required field 'items'")
	}
	if len(order.Products) == 0 {
		t.Error("Order: missing required field 'products'")
	}

	for i, item := range order.Items {
		if item.ProductID == "" {
			t.Errorf("Order item[%d]: missing required field 'productId'", i)
		}
		if item.Quantity <= 0 {
			t.Errorf("Order item[%d]: invalid quantity value: %d", i, item.Quantity)
		}
	}

	for i, product := range order.Products {
		validateProductSchema(t, product, fmt.Sprintf("order.products[%d]", i))
	}
}

func validateErrorResponse(t *testing.T, resp *http.Response) {
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Error response: Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var errorResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	errorField, exists := errorResp["error"]
	if !exists {
		t.Error("Error response: missing 'error' field")
		return
	}

	errorObj, ok := errorField.(map[string]interface{})
	if !ok {
		t.Error("Error response: 'error' field is not an object")
		return
	}

	if _, exists := errorObj["code"]; !exists {
		t.Error("Error response: missing 'error.code' field")
	}
	if _, exists := errorObj["type"]; !exists {
		t.Error("Error response: missing 'error.type' field")
	}
	if _, exists := errorObj["message"]; !exists {
		t.Error("Error response: missing 'error.message' field")
	}
}

func TestResponseTimeCompliance(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.server.Close()

	t.Run("Product listing response time", func(t *testing.T) {
		start := time.Now()
		resp, err := http.Get(testServer.server.URL + "/product")
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if duration > 500*time.Millisecond {
			t.Errorf("Product listing took %v, expected under 500ms", duration)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
