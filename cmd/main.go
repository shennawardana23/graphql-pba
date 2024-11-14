package main

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/shennawardana23/graphql-pba/graph"
	"github.com/shennawardana23/graphql-pba/graph/generated"
	"github.com/shennawardana23/graphql-pba/internal/app/database"
	"github.com/shennawardana23/graphql-pba/internal/repository"
	"github.com/shennawardana23/graphql-pba/internal/service"
	"github.com/shennawardana23/graphql-pba/internal/util/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
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
	db := database.NewDB()
	validate := validator.New()
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo, db, validate)

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			UserService: userService,
		},
	}))
	log.Println("GraphQL server created successfully")

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add custom logging middleware
	r.Use(gin.Recovery())
	r.Use(loggerMiddleware())

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// GraphQL endpoints
	r.POST("/query", gin.WrapH(srv))
	r.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	port := os.Getenv("PORT")
	// Log startup time
	logger.Log.WithFields(logrus.Fields{
		"startup_time": time.Since(startTime).String(),
		"port":         port,
	}).Info("Server is starting up on port " + port)

	// Start server
	log.Printf("GraphQL playground available at http://localhost:%s/", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

// Custom logging middleware for Gin
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Increment the request count
		requestCount.WithLabelValues(c.Request.Method).Inc() // Increment request count

		// Log to file
		log.Printf("method=%s path=%s status=%d latency=%s client_ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start).String(),
			c.ClientIP(),
		)
	}
}
