package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"flippOneApi/internal/api"
)

type Application struct {
	Logger      *log.Logger
	UserHandler *api.UserHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		Logger:      logger,
		UserHandler: api.NewUserHandler(),
	}

	return app, nil
}

func (app *Application) HealthCheck(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("Health check %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(res, "OK\n")
}
