package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ooliokartchallenge/internal/application/services"
	"ooliokartchallenge/internal/config"
	"ooliokartchallenge/internal/domain/interfaces"
	httpInfra "ooliokartchallenge/internal/infrastruture/http"
	"ooliokartchallenge/internal/infrastruture/http/handlers"
	"ooliokartchallenge/internal/infrastruture/http/middleware"
	"ooliokartchallenge/internal/infrastruture/repositories"
	"ooliokartchallenge/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config         *config.Config
	logger         *logger.Logger
	server         *http.Server
	productRepo    interfaces.ProductRepository
	promoRepo      interfaces.PromoRepository
	productSerivce interfaces.ProductService
	promoService   interfaces.PromoService
	orderService   interfaces.OrderService
}

func main() {

	app, err := initialzeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v ", err)

	}

	if err := app.start(); err != nil {
		app.logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}

func initialzeApp() (*App, error) {
	cfg := config.Load()

	appLogger := logger.New()
	appLogger.Info("Configuration loaded successfully")

	productRepo := repositories.NewProductRepository()

	appLogger.Info("Initializing promo repository", "files", cfg.CouponFiles)
	promoRepo := repositories.NewPromoRepository(cfg.CouponFiles)



	appLogger.Info("Initializing application services")

	promoService := services.NewPromoService(promoRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(productRepo, promoService)

	ctx := context.Background()

	if _, err := productService.ListProducts(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize product service: %w", err)
	}

	appLogger.Info("All services initialized successfully")

	productHandler := handlers.NewProductHandler(productService, appLogger)
	orderHandler := handlers.NewOrderHandler(orderService, appLogger)

	authMiddlerware := middleware.NewAuthMiddleware(appLogger)
	corsMiddleware := middleware.NewCORSMiddleware()

	router := httpInfra.NewRouter(productHandler, orderHandler, authMiddlerware, corsMiddleware)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.SetupRoutes(),
	}

	return &App{
		config:         cfg,
		logger:         appLogger,
		server:         server,
		productRepo:    productRepo,
		promoRepo:      promoRepo,
		productSerivce: productService,
		promoService:   promoService,
		orderService:   orderService,
	}, nil

}

func (a *App) start() error {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info("Starting HTTP server", "address", a.server.Addr)

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Server failed to start", "error", err)
			quit <- syscall.SIGTERM

		}
	}()

	a.logger.Info("Server is running successfully",
		"port", a.config.Port,
		"api_key_configured", a.config.APIKey != "",
		"promo_file_loaded", len(a.config.CouponFiles))

	<-quit

	a.logger.Info("Sutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", "error", err)
	}

	a.logger.Info("Server shutdown complete")

	return nil

}
