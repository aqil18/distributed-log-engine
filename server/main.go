package main

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got / request\n")
    io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /hello request\n")
    io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
    http.HandleFunc("/", getRoot)
    http.HandleFunc("/hello", getHello)
    
    err := http.ListenAndServe(":3333", nil)

}

// ALTERNATIVE
type helloHandler struct {
    db *sql.DB
}

// receiver that takes in a state
func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from a handler with a DB connection!")
}
// we can then use our handler in listen and serve as it counts as a Handler
http.ListenAndServe(":8080", &helloHandler{db: myDB})
