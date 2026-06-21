package main



import (
    "context"
    "errors"
    "fmt"
    "io"
    "net"
    "net/http"
)

const keyServerAddr = "serverAddr"


func getRoot(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // key value lookup of serverAddr in context    
    fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
    io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /hello request\n")
    io.WriteString(w, "Hello, HTTP!\n")
}





func readEntry(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /readEntry request\n")
    io.WriteString(w, "Reading entry!\n")

    body, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("could not read body: %s\n", err)
    }

    fmt.Printf("%s: got / request. body:\n%s\n",
        body)
}



// func appendEntry(w http.ResponseWriter, r *http.Request) {
//     fmt.Printf("got /hello request\n")
//     io.WriteString(w, "Hello, HTTP!\n")

//     // Parse request to find entry and call appendEntry
//     LogEntry entry
//     appendEntry(entry)
// }

func main() {
    mux := http.NewServeMux() 
	// using an explicit server multiplexer 
	// default one can quickly be populated with other servers
    mux.HandleFunc("/", getRoot)
    mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/readEntry" , readEntry)
	// mux.HandleFunc("/appendentry" , appendEntry)


    ctx, cancelCtx := context.WithCancel(context.Background())
    serverOne := &http.Server{
        Addr:    ":3333",
        Handler: mux,
        BaseContext: func(l net.Listener) context.Context {
            // create a context and add the keyServerAddr and address to it
            ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
            return ctx
        },
    }


    serverTwo := &http.Server{
        Addr:    ":4444",
        Handler: mux,
        BaseContext: func(l net.Listener) context.Context {
            ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
            return ctx
        },
    }

    // Kicks off a go routine for server 1 and 2 and they run concurrently
    go func() {
        err := serverOne.ListenAndServe()
        if errors.Is(err, http.ErrServerClosed) {
            fmt.Printf("server one closed\n")
        } else if err != nil {
            fmt.Printf("error listening for server one: %s\n", err)
        }
    
        cancelCtx()
    }()

    go func() {
    err := serverTwo.ListenAndServe()
    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server two closed\n")
    } else if err != nil {
        fmt.Printf("error listening for server two: %s\n", err)
    }
        cancelCtx()
    }()
    
    // This is to do with channel management - look into this
    <-ctx.Done()
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
