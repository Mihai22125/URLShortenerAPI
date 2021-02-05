package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mihai22125/URLShortenerAPI/data"
	"github.com/Mihai22125/URLShortenerAPI/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	_ "github.com/pdrum/swagger-automation/docs" // This line is necessary for go-swagger to find docs!
)

var bindAddress = env.String("BIND_ADRESS", false, ":8080", "Bind address for the server")

func main() {

	env.Parse()

	l := log.New(os.Stdout, "url-api", log.LstdFlags)
	urlList := data.Urls{}

	// create the handlers
	uh := handlers.NewUrls(l, urlList)

	// create a new serve mux
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/v1/{shortURL}", uh.GetURL)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/v1/", uh.AddURL)
	postRouter.Use(uh.MiddlewareValidateURL)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// create a new server

	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connecting using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port ", s.Addr)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-c
	log.Println("Got signal: ", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		l.Println("Server Shutdown Failed", "error", err)
	} else {
		l.Println("Server Shutdown gracefully")
	}
}
