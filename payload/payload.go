package payload

import (
  "io"
  "fmt"
  "net/http"
  "os"
  "github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)

    io.WriteString(w, "Crunchy Data PostgreSQL for Pivotal Cloud Foundry Upload Utility\n")
    io.WriteString(w, "Data File Upload\n")
    fmt.Fprintf(w, "Utility: %v\n", vars["utility"])

    file, handler, err := r.FormFile("file")
    if err != nil {
        panic(err) //dont do this
    }
    defer file.Close()

    // copy example
    f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        panic(err) //please dont
    }
    defer f.Close()
    io.Copy(f, file)

}
