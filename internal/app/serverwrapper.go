package server

import (
	sqlstore "challenge/internal/store/sqlstore"
	router "challenge/pkg"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"gitlab.com/0x4149/logz"
)

const (
	GITHUB_AT = "GITHUB_AT"
)

func Start(config *Config) error {
	//create database
	db, err := newDB(config.DatabaseURL, config.DatabaseSchema)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)

	//create CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	//create router
	router := router.New()

	//new http.server with corsmiddleware
	newHttpSrv := &http.Server{
		Addr:    ":8020",
		Handler: corsMiddleware.Handler(router),
	}

	//create github app

	loadEnv()
	//Starting github app to scout users
	app := NewGithubAPP(os.Getenv(GITHUB_AT), store)

	//Adding all users i want to scout
	// app.AddUser("flight505")

	//starting github
	go app.Start()

	srv := newServer(newHttpSrv, app, store, router)

	logz.Info("Starting a server on port", srv.Addr)
	return srv.Server.ListenAndServe()
}

func newDB(databaseURL, dataBaseSchema string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, err
	}

	sqlStmt, err := os.ReadFile(dataBaseSchema)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(string(sqlStmt))
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		logz.Error("Error loading .env file")
		os.Exit(1)
	}
}
