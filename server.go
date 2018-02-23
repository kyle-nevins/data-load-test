package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "io"
    "os"
    "fmt"
    "test-load/payload"
)

// Redirect any requests to the base page to the Crunchy Data main page.
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "http://www.crunchydata.com", 303)
}

func main() {
    r := mux.NewRouter()

    dataPayload := payload.NewPayload()

    r.HandleFunc("/upload/{utility}", *dataPayload.UploadHandler)

    // Catch All
    //r.PathPrefix("/").HandleFunc(catchAllHandler)

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8000", r))
}
