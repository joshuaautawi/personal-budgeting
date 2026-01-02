package main

import (
	"log"
	"os"

	"personal-budgeting/be/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	a := app.New()
	log.Printf("listening on :%s", port)
	log.Fatal(a.Listen(":" + port))
}
