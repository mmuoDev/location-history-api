package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mmuoDev/location-history-api.git/internal/app"
	"github.com/mmuoDev/location-history-api.git/pkg/db"
)

func main() {
	port := os.Getenv("HISTORY_SERVER_LISTEN_ADDR")
	conn := db.New()
	a := app.New(conn)
	if port == "" {
		port = "8080"
	}
	log.Println(fmt.Sprintf("Starting server on port:%s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), a.Handler()))
}
