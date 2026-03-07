package main

import (
	"fmt"
	appAuth "graduation-invitation/internal/app/auth"
	appConfig "graduation-invitation/internal/app/config"
	"graduation-invitation/internal/app/guest"
	"graduation-invitation/internal/config"
	"graduation-invitation/internal/infra/database"
	"graduation-invitation/internal/infra/persistence"
	"graduation-invitation/internal/infra/transport/http"
	stdhttp "net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := database.InitPostgres(dsn)
	if err != nil {
		panic(err)
	}

	// Auth
	authSvc := appAuth.NewAuthService(cfg.AdminUsername, cfg.AdminPasswordHash, cfg.JWTSecret)
	authHandler := http.NewAuthHandler(authSvc)

	// Guest
	guestRepo := persistence.NewGuestPersistence(db)
	guestSvc := guest.NewGuestService(guestRepo)
	guestHandler := http.NewGuestHandler(guestSvc)

	// Config
	configRepo := persistence.NewConfigPersistence(db)
	configSvc := appConfig.NewConfigService(configRepo)
	configHandler := http.NewConfigHandler(configSvc)

	// Upload
	uploadHandler := http.NewUploadHandler(cfg.UploadDir)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORSAllowed,
		AllowMethods: []string{stdhttp.MethodGet, stdhttp.MethodPost, stdhttp.MethodPut, stdhttp.MethodDelete, stdhttp.MethodPatch, stdhttp.MethodOptions, stdhttp.MethodHead},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Routes
	http.RegisterRoutes(e, guestHandler, configHandler, authHandler, uploadHandler, cfg.JWTSecret, cfg.UploadDir)

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
