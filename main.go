package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/leonardosantos2/GoUserServer/internal/app"
	"github.com/leonardosantos2/GoUserServer/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "Port running server")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	router := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("Starting server on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
