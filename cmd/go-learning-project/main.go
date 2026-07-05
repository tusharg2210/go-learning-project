package main

import (
	"fmt"
	"log" // standard logging
	"log/slog" // structured logging
	"net/http" // standard HTTP server
	"os" // for OS-level operations
	"os/signal" // for handling OS signals
	"syscall" // for handling system calls
	"context" // for context management
	"time" // for time-related functions
	"github.com/tusharg2210/go-learning-project/internal/config"
)

func main() {
	// Load the configuration using MustLoadConfig
	cfg := config.MustLoadConfig()


	// setup logger

	// database setup

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r * http.Request){
		w.Write([]byte("Welcome to go"))
	})
	// setup server
	server := http.Server {
		Addr:    cfg.Address,
		Handler: router,
	}
	
	slog.Info("Starting server", slog.String("address", cfg.Address))
	fmt.Printf("Server started on %s", cfg.Address)

	// For gracefully shutdown the server, we can use a go routine to listen for OS signals and call server.Shutdown() when a termination signal is received. This will allow the server to finish processing any ongoing requests before shutting down.
	done := make(chan os.Signal, 1) // buffered channel to receive OS signals
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // listen for termination signals and notify the done channel

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	} ()
	
	<- done // wait for termination signal

	slog.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // create a context with a timeout for graceful shutdown
	defer cancel() // ensure the cancel function is called to release resources

	
	// err := server.Shutdown(ctx) // gracefully shutdown the server
	// if err != nil {
		// 	slog.Error("failed to shutdown sever", slog.String("error", err.Error()))
		// }
		
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown sever", slog.String("error", err.Error()))
	}
	
	slog.Info("Server stopped successfully")

}
