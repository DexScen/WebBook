package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	msql "github.com/DexScen/WebBook/backend/auth/internal/repository/psql"
	"github.com/DexScen/WebBook/backend/auth/internal/service"
	"github.com/DexScen/WebBook/backend/auth/internal/transport/rest"
	"github.com/DexScen/WebBook/backend/auth/pkg/database"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  "disable",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usersRepo := msql.NewUsers(db)
	usersService := service.NewBooks(usersRepo)
	handler := rest.NewHandler(usersService)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: handler.InitRouter(),
	}

	log.Println("Server started at:", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
