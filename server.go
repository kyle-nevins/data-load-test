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
    "strings"
    "path/filepath"
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
  h := md5.New()
  token := fmt.Sprintf("%x", h.Sum(nil))
  templ := template.New("upload")

  if r.Method == "GET" {
    crutime := time.Now().Unix()
    io.WriteString(h, strconv.FormatInt(crutime, 10))

    switch vars["utility"]{
      case "pgloader":
        if _, err := templ.ParseFiles("web/pgloader.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
        if err := templ.ExecuteTemplate(w, "pgloader", token); err != nil { panic(err) }
      case "psql":
        if _, err := templ.ParseFiles("web/psql.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
        if err := templ.ExecuteTemplate(w, "psql", token); err != nil { panic(err) }
      case "pgrestore":
        if _, err := templ.ParseFiles("web/pgrestore.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
        if err := templ.ExecuteTemplate(w, "pgrestore", token); err != nil { panic(err) }
    }

  } else {
    dataPayload := new(payload.Payload)

    if validFile := ValidateFile(w,r,"data"); validFile {
      dataPayload.AssignVariables(w,r)
      dataPayload.WriteFile(w,r)
      dataPayload.DataLoader(w,r)
    } else {
      if _, err := templ.ParseFiles("web/error.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
      if err := templ.ExecuteTemplate(w, "error", token); err != nil { panic(err) }
    }
  }
}

func ValidateFile(w http.ResponseWriter, r *http.Request, formField string) bool {

  validFile := false

  file, header, err := r.FormFile(formField)
  switch err {
    case nil:
      log.Println("File: ",header.Filename)
      checkBuf := make([]byte, 512)
    	if _, err := io.ReadAtLeast(file, checkBuf, 512); err != nil {
    		log.Println("Check Buffer Error:", err)
    	}
      allowedExtensions := []string{".csv",".sql",".dump"}
      fileExt := filepath.Ext(header.Filename)
      for _, extensions := range allowedExtensions {
        if extensions == fileExt {
          if filetype := http.DetectContentType(checkBuf); (strings.Contains(filetype, "text/plain")) {
            log.Println("File Type:", filetype)
            validFile = true
          } else {
            log.Println("Not a text file:", filetype)
          }
        }
      }
    case http.ErrMissingFile:
      log.Println("no file")
    default:
      log.Println("File Read Error:", err)
  }

  if ( mux.Vars(r)["utility"] == "pgloader" ) {
    validFile = ValidateFile(w,r,"loaderFile")
  }

  return validFile
}

// Redirect any requests to the base page to the Crunchy Data main page.
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "http://www.crunchydata.com", 303)
}
