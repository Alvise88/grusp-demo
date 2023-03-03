package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"grusp.io/hello/cmd/web/handlers"
)

func main() {
	ctx := context.Background()

	port := "8080"

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           handlers.Routes(ctx),
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("unable to run argoci cli")
	}
}
