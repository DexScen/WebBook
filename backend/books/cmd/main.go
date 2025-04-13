package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DexScen/WebBook/backend/books/internal/repository/psql"
	"github.com/DexScen/WebBook/backend/books/internal/service"
	"github.com/DexScen/WebBook/backend/books/internal/transport/rest"
	"github.com/DexScen/WebBook/backend/books/pkg/database"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode: "disable",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	booksRepo := msql.NewBooks(db)
	booksService := service.NewBooks(booksRepo)
	handler := rest.NewHandler(booksService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("Server started at:", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
