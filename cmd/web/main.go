package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"volchok96.com/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	snippets    *mysql.SnippetModel
	redisClient *redis.Client
}

func connectToRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // Address of the Redis server
	})

	// Check the connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func main() {
	addr := flag.String("addr", ":4000", "Network address of the web server")
	dsn := flag.String("dsn", "/snippetbox?parseTime=true", "MySQL data source name")
	redisAddr := flag.String("redis", "localhost:6379", "Redis server address") // Added for Redis address
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Prompt for the database password
	password, err := promptForPassword()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Create the full DSN with the entered password
	fullDSN := constructDSN(password, *dsn)

	// Create a connection pool to the database
	db, err := openDB(fullDSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Connect to Redis
	rdb, err := connectToRedis(*redisAddr)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer rdb.Close()

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		snippets:    &mysql.SnippetModel{DB: db},
		redisClient: rdb, // Storing the Redis client in the application structure
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
