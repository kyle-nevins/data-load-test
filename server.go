package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "data-load-test/payload"
    "io"
    "html/template"
)

func main() {
    r := mux.NewRouter()


    // Handler function for uploading data files.
    r.HandleFunc("/upload/{utility}", UploadHandler)

    // Handler function for displaying Open Layers
    r.HandleFunc("/display/{table}", DisplayHandler)

    // Catch All
    //r.PathPrefix("/").HandleFunc(catchAllHandler)

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8080", r))
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

  dataPayload := new(payload.Payload)

  dataPayload.AssignVariables(w,r)
  dataPayload.WriteFile(w,r)
  dataPayload.DataLoader(w,r)

  io.WriteString(w, "Data Loader run completed.\n")
}

type MapPageData struct {
  PageTitle string
}

func DisplayHandler(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "Display Handler")
  tmpl := template.Must(template.ParseFiles("web/map.html"))
  data := MapPageData{
			PageTitle: "Crunchy PostgreSQL PG Data Load Map Display",
		}
		tmpl.Execute(w, data)
}

// Redirect any requests to the base page to the Crunchy Data main page.
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "http://www.crunchydata.com", 303)
}
