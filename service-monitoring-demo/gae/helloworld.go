package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10) // n will be between 0 and 10
		fmt.Printf("randon number was %d\n", n)

		fmt.Fprintf(w, "Hello World!")

	})

	http.ListenAndServe(":8080", r)
}
