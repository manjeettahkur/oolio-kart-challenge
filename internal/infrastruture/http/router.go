package http

import (
	"net/http"
	"ooliokartchallenge/internal/infrastruture/http/handlers"
	"ooliokartchallenge/internal/infrastruture/http/middleware"
)

type Router struct {
	productHandler *handlers.ProductHandler
	orderHandler   *handlers.OrderHandler
	authMiddleware *middleware.AuthMiddleware
	corsMiddleware *middleware.CORSMiddleware
}

func NewRouter(
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	authMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CORSMiddleware,
) *Router {
	return &Router{
		productHandler: productHandler,
		orderHandler:   orderHandler,
		authMiddleware: authMiddleware,
		corsMiddleware: corsMiddleware,
	}
}

// SetupRoutes configures the application's HTTP routes using the enhanced Go 1.22+ ServeMux.
func (r *Router) SetupRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /product", r.productHandler.ListProducts)
	mux.HandleFunc("GET /product/{id}", r.productHandler.GetProduct)

	protectedOrderHandler := r.authMiddleware.RequireAPIKey(http.HandlerFunc(r.orderHandler.PlaceOrder))
	mux.Handle("POST /order", protectedOrderHandler)

	finalHandler := r.corsMiddleware.EnableCORS(mux)

	return finalHandler
}
