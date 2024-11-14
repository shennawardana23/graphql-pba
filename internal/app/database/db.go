package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/shennawardana23/graphql-pba/internal/util/logger"
	"github.com/sirupsen/logrus"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbName   = os.Getenv("DB_NAME")
)

// func Connect() *pg.DB {

// 	// Use pg.Options to set up the connection parameters
// 	opt := &pg.Options{
// 		Addr:     fmt.Sprintf("%s:%s", "localhost", "5432"),
// 		User:     "postgres",
// 		Password: "postgres",
// 		Database: "auth_db",
// 		// SSLMode:  "disable",
// 	}

// 	// Connect to the database
// 	db := pg.Connect(opt)

// 	// Check if the database is reachable
// 	if _, err := db.Exec("SELECT 1"); err != nil {
// 		panic("PostgreSQL is down: " + err.Error())
// 	}

// 	return db
// }

// query interface
type Q interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB interface {
	Q
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(ctx context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Ping() error
	PingContext(ctx context.Context) error
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}

type TX interface {
	Q
	Commit() error
	Rollback() error
	Stmt(stmt *sql.Stmt) *sql.Stmt
	StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt
}

type store struct {
	db *sql.DB
}

func (s *store) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}
func (s *store) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}
func (s *store) Close() error {
	return s.db.Close()
}
func (s *store) Conn(ctx context.Context) (*sql.Conn, error) {
	return s.db.Conn(ctx)
}
func (s *store) Driver() driver.Driver {
	return s.db.Driver()
}
func (s *store) Exec(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args...)
}
func (s *store) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
func (s *store) Ping() error {
	return s.db.Ping()
}
func (s *store) PingContext(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
func (s *store) Prepare(query string) (*sql.Stmt, error) {
	return s.db.Prepare(query)
}
func (s *store) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, query)
}
func (s *store) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}
func (s *store) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}
func (s *store) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args...)
}
func (s *store) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}
func (s *store) SetConnMaxIdleTime(d time.Duration) {
	s.db.SetConnMaxIdleTime(d)
}
func (s *store) SetConnMaxLifetime(d time.Duration) {
	s.db.SetConnMaxLifetime(d)
}
func (s *store) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}
func (s *store) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}
func (s *store) Stats() sql.DBStats {
	return s.db.Stats()
}

func NewDB() *sql.DB {
	return newDb(dbName)
}

func NewTestDb() *sql.DB {
	return newDb(dbName)
}

func newDb(dbName string) *sql.DB {

	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbName,
	)

	logger.WithFields(logrus.Fields{
		"host":              host,
		"port":              port,
		"db":                dbName,
		"user":              user,
		"connection_string": postgresInfo,
	}).Info("Successfully established a connection into Postgres DB!")

	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to open database connection")
	}
	ctx := context.TODO()
	if err = db.Ping(); err != nil {
		if err = db.Close(); err != nil {
			fmt.Print("Ping failed into PostgresDB !", err, ctx)
			logger.Warn(ctx, err)
		}
	}
	fmt.Println("Established a successful connection into Postgres DB!")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(time.Minute * 5)

	return db
}
