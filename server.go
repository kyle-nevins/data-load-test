package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "data-load-test/payload"
    "io"
    "html/template"
    "fmt"
    "strconv"
    "crypto/md5"
    "time"
)

func main() {
    r := mux.NewRouter()

    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

    // Handler function for uploading data files.
    r.HandleFunc("/cli/{utility}", CLIHandler)

    // Handler function for displaying Open Layers
    r.HandleFunc("/web/{utility}", WebHandler)

    // Catch All
    //r.PathPrefix("/").HandleFunc(catchAllHandler)

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8080", r))
}

func CLIHandler(w http.ResponseWriter, r *http.Request) {

  dataPayload := new(payload.Payload)

  dataPayload.AssignVariables(w,r)
  dataPayload.WriteFile(w,r)
  dataPayload.DataLoader(w,r)

  io.WriteString(w, "Data Loader run completed.\n")
}

func WebHandler(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)

  if r.Method == "GET" {
    crutime := time.Now().Unix()
    h := md5.New()
    io.WriteString(h, strconv.FormatInt(crutime, 10))
    token := fmt.Sprintf("%x", h.Sum(nil))

    t := template.New("upload")

    switch vars["utility"]{
      case "pgloader":
        if _, err := t.ParseFiles("web/pgloader.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
        if err := t.ExecuteTemplate(w, "pgloader", token); err != nil { panic(err) }
      case "psql":
        if _, err := t.ParseFiles("web/psql.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
        if err := t.ExecuteTemplate(w, "psql", token); err != nil { panic(err) }
      case "pgrestore":
        if _, err := t.ParseFiles("web/pgrestore.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
        if err := t.ExecuteTemplate(w, "pgrestore", token); err != nil { panic(err) }
    }

  } else {
    dataPayload := new(payload.Payload)
    dataPayload.AssignVariables(w,r)
    dataPayload.WriteFile(w,r)
    dataPayload.DataLoader(w,r)
  }
}

// Redirect any requests to the base page to the Crunchy Data main page.
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "http://www.crunchydata.com", 303)
}
