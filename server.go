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
    "os"
)

func main() {
    r := mux.NewRouter()
    log.SetPrefix("CrunchyDataDataLoader: ")

    r.HandleFunc("/favicon.ico", faviconHandler)
    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(os.Getenv("HOME")+"/static/css"))))
    r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(os.Getenv("HOME")+"/static/images"))))
    r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(os.Getenv("HOME")+"/static/js"))))


    // Handler function for uploading data files.
    r.HandleFunc("/cli/{utility}", CLIHandler)

    // Handler function for displaying Open Layers
    r.HandleFunc("/web/{utility}", WebHandler)

    // Handler function for default webpage.
    r.HandleFunc("/",HomePage)

    //Pass in a handler for any 404 pages.
    r.NotFoundHandler = http.HandlerFunc(notFound)

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8080", r))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/images/crunchy-left.png")
}

func notFound(w http.ResponseWriter, r *http.Request) {
  log.Println("URL Not Found:", r.URL.Path)
  http.ServeFile(w, r, "web/notfound.html")
}

func HomePage(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "web/index.html")
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
  crutime := time.Now().Unix()
  io.WriteString(h, strconv.FormatInt(crutime, 10))
  token := fmt.Sprintf("%x", h.Sum(nil))
  t := template.New("upload")

  if r.Method == "GET" {
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
      default:
        http.ServeFile(w, r, "web/notfound.html")
    }
  } else {
    dataPayload := new(payload.Payload)

    if validFile := ValidateFile(w,r,"data"); validFile {
      dataPayload.AssignVariables(w,r)
      dataPayload.WriteFile(w,r)
      dataPayload.DataLoader(w,r)
    } else {
      if _, err := t.ParseFiles("web/error.gtpl", "web/layout.gtpl"); err != nil { panic(err) }
      if err := t.ExecuteTemplate(w, "error", token); err != nil { panic(err) }
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
          log.Printf("File Extension Found: %s", extensions)
          if filetype := http.DetectContentType(checkBuf); (strings.Contains(filetype, "text/plain")) {
            log.Println("File Type:", filetype)
            validFile = true
          } else {
            log.Println("Not a text file:", filetype)
          }
        }
      }
      if !validFile {
        log.Println("No Valid File Extension.")
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
