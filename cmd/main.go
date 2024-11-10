package main

import (
	"context"
	"log"
	"net/http"
	"notion-connect/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notion-connect/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Use configuration values
	// serverAddr := cfg.Server.Host + ":" + cfg.Server.Port

	// // Initialize database connection
	// db, err := initDB(cfg.Database.GetDSN())
	// if err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }

	// // Initialize Redis if needed
	// redis, err := initRedis(cfg.Redis)
	// if err != nil {
	//     log.Fatalf("Failed to connect to Redis: %v", err)
	// }

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set Gin mode
	// gin.SetMode(GetEnv("GIN_MODE"))
	if cfg.IsDevelopment() {
		gin.SetMode("debug")
	}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	// Define a simple route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	notion := middleware.NotionConnectInit(cfg.NotionAPI)
	router.GET("/notion/database", notion.ConnectNotion)
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server
	startServer(srv)

}

func startServer(srv *http.Server) {
	// Graceful shutdown
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
