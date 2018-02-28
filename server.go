package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "data-load-test/payload"
    "io"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/upload/{utility}", UploadHandler)

    dataPayload := new(payload.Payload)

    r.Use(dataPayload.AssignVariables)
    r.Use(dataPayload.WriteFile)
    r.Use(dataPayload.DataLoader)

    // Catch All
    //r.PathPrefix("/").HandleFunc(catchAllHandler)

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8080", r))
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "Data Loader run completed.\n")
}

// Redirect any requests to the base page to the Crunchy Data main page.
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "http://www.crunchydata.com", 303)
}
