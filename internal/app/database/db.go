package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/shennawardana23/graphql-pba/internal/util/logger"
)

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	PoolSize        int
	MinIdleConns    int
	MaxConnAge      time.Duration
	PoolTimeout     time.Duration
	IdleTimeout     time.Duration
	MaxRetries      int
	MaxRetryBackoff time.Duration
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:            getEnvOrDefault("DB_HOST", "localhost"),
		Port:            getEnvOrDefault("DB_PORT", "5432"),
		User:            getEnvOrDefault("DB_USER", "postgres"),
		Password:        getEnvOrDefault("DB_PASSWORD", "postgres"),
		Database:        getEnvOrDefault("DB_NAME", "auth_db"),
		PoolSize:        getEnvAsInt("DB_POOL_SIZE", 10),
		MinIdleConns:    getEnvAsInt("DB_MIN_IDLE_CONNS", 5),
		MaxConnAge:      getEnvAsDuration("DB_MAX_CONN_AGE", "1h"),
		PoolTimeout:     getEnvAsDuration("DB_POOL_TIMEOUT", "30s"),
		IdleTimeout:     getEnvAsDuration("DB_IDLE_TIMEOUT", "5m"),
		MaxRetries:      getEnvAsInt("DB_MAX_RETRIES", 3),
		MaxRetryBackoff: getEnvAsDuration("DB_MAX_RETRY_BACKOFF", "5s"),
	}
}

func Connect() *pg.DB {
	config := NewDBConfig()

	opt := &pg.Options{
		Addr:            fmt.Sprintf("%s:%s", config.Host, config.Port),
		User:            config.User,
		Password:        config.Password,
		Database:        config.Database,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxConnAge:      config.MaxConnAge,
		PoolTimeout:     config.PoolTimeout,
		IdleTimeout:     config.IdleTimeout,
		MaxRetries:      config.MaxRetries,
		MaxRetryBackoff: config.MaxRetryBackoff,

		// Enable logging for slow queries
		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			logger.Log.Info("New database connection established")
			return nil
		},
	}

	db := pg.Connect(opt)

	// Add hooks for query monitoring
	db.AddQueryHook(queryHook{})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Log database stats periodically
	go monitorDBStats(db)

	return db
}

// QueryHook for monitoring queries
type queryHook struct{}

func (h queryHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h queryHook) AfterQuery(ctx context.Context, evt *pg.QueryEvent) error {
	query, err := evt.FormattedQuery()
	if err != nil {
		return err
	}

	logger.Log.WithFields(map[string]interface{}{
		"query":    string(query),
		"duration": time.Since(evt.StartTime),
	}).Debug("Database query executed")

	return nil
}

// Monitor database stats
func monitorDBStats(db *pg.DB) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := db.PoolStats()
		logger.Log.WithFields(map[string]interface{}{
			"total_conns": stats.TotalConns,
			"idle_conns":  stats.IdleConns,
			"stale_conns": stats.StaleConns,
			"hits":        stats.Hits,
			"misses":      stats.Misses,
			"timeouts":    stats.Timeouts,
		}).Info("Database connection pool stats")
	}
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	value := getEnvOrDefault(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		logger.Log.Warnf("Invalid duration for %s, using default: %s", key, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}
