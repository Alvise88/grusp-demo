package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"grusp.io/hello/cmd/web/application"
	"grusp.io/hello/cmd/web/handlers"
	"grusp.io/hello/cmd/web/templates"
)

func main() {
	ctx := context.Background()

	port := "8080"

	templateCache, err := templates.NewTemplateCache()

	if err != nil {
		panic(err)
	}

	app := application.NewApplication(templateCache)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           handlers.Routes(ctx, app),
	}

	err = server.ListenAndServe()

	if err != nil {
		fmt.Println("unable to run argoci cli")
	}
}
