package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shennawardana23/graphql-pba/graph"
	"github.com/shennawardana23/graphql-pba/graph/generated"
	"github.com/shennawardana23/graphql-pba/internal/app/database"
	"github.com/shennawardana23/graphql-pba/internal/middleware"
	"github.com/shennawardana23/graphql-pba/internal/util/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

const (
	shutdownTimeout = 5 * time.Second
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_requests_total",
			Help: "Total number of GraphQL requests",
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
}

func main() {
	startTime := time.Now()

	// Setup signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a logs directory if it doesn't exist
	os.MkdirAll("logs", os.ModePerm)

	// Create a log file
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set log output to the log file
	log.SetOutput(logFile)

	// Initialize database
	db := database.Connect()
	defer func() {
		logger.Log.Info("Closing database connection...")
		if err := db.Close(); err != nil {
			logger.Log.Errorf("Error closing database connection: %v", err)
		}
	}()

	// Create resolver with dependencies
	resolver := graph.NewResolver(db)

	// Create GraphQL server with custom error presenter
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Set custom error presenter
	srv.SetErrorPresenter(graph.ErrorPresenter)

	log.Println("GraphQL server created successfully")

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add custom logging middleware
	r.Use(gin.Recovery())
	r.Use(loggerMiddleware())
	r.Use(middleware.ErrorHandler())

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// GraphQL endpoints
	r.POST("/query", gin.WrapH(srv))
	r.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	port := os.Getenv("PORT")

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Log.WithFields(logrus.Fields{
			"startup_time": time.Since(startTime).String(),
			"port":         port,
		}).Info("Server is starting up on port " + port)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Log.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown tasks
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
	}

	// Wait for context to complete
	<-ctx.Done()
	logger.Log.Info("Server shutdown complete")
}

// Custom logging middleware for Gin
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Add request tracking
		c.Next()

		// Skip metrics endpoint logging
		if path != "/metrics" {
			// Increment the request count
			requestCount.WithLabelValues(c.Request.Method).Inc()

			// Log request details
			logger.Log.WithFields(logrus.Fields{
				"method":     c.Request.Method,
				"path":       path,
				"status":     c.Writer.Status(),
				"latency":    time.Since(start).String(),
				"client":     c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Info("Request processed")
		}
	}
}
