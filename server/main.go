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



func appendEntry(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /hello request\n")
    io.WriteString(w, "Hello, HTTP!\n")

    // Parse request to find entry and call appendEntry
    LogEntry entry
    appendEntry(entry)
}

func readEntry(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /hello request\n")
    io.WriteString(w, "Hello, HTTP!\n")

    // Parse request to find entry index and call readEntryByIndex
    
    readEntryByIndex(0)
    LogEntry entry
}


func main() {
    mux := http.NewServeMux() 
	// using an explicit server multiplexer 
	// default one can quickly be populated with other servers
    mux.HandleFunc("/", getRoot)
    mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/appendentry" , appendEntry)
	mux.HandleFunc("/readentry" , readentry)


    err := http.ListenAndServe(":3333", mux)

	if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed\n")
    } else if err != nil {
        fmt.Printf("error starting server: %s\n", err)
        os.Exit(1)
    }

}

// ALTERNATIVE TO ADD STATE
// type helloHandler struct {
//     db *sql.DB
// }

// // receiver that takes in a state
// func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "Hello from a handler with a DB connection!")
// }
// // we can then use our handler in listen and serve as it counts as a Handler
// http.ListenAndServe(":8080", &helloHandler{db: myDB})
