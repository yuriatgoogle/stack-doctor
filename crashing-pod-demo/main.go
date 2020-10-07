package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	// Create a mux for routing incoming requests
	m := http.NewServeMux()

	// All URLs will be handled by this function
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("serving request")
		w.Write([]byte("Hello, world!"))
	})

	// crash the process
	m.HandleFunc("/crashme", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("crash endpoint called - exiting!")
		os.Exit(1)
	})

	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	// Continue to process new requests until an error occurs
	log.Fatal(s.ListenAndServe())
}
