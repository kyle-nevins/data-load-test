package payload

import (
  "io"
  "fmt"
  "net/http"
  "os"
  "github.com/gorilla/mux"
  //"strings"
)

type Payload struct {
  payloadFile string
  Utility  string
}

// Middleware function, which will be called for each request
func (dp *Payload) AssignVariables(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dp.Utility = vars["utility"]
    fmt.Fprintf(w, "Upload Utility %v\n", dp.Utility)
    next.ServeHTTP(w, r)
  })
}

func (dp *Payload) WriteFile(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    switch dp.Utility {
      case "pgloader":
        loaderFile, loaderHandler, loaderErr := r.FormFile("loaderFile")
        if loaderErr != nil { panic(loaderErr) }
        dataFile, dataHandler, dataErr := r.FormFile("data")
        if dataErr != nil { panic(dataErr) }
        defer loaderFile.Close()
        defer dataFile.Close()
        lf, le := os.OpenFile(os.Getenv("HOME")+"/pgloader/loaderFile", os.O_WRONLY|os.O_CREATE, 0666)
        if le != nil { panic(le) }
        df, de := os.OpenFile(os.Getenv("HOME")+"/pgloader/"+dataHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if de != nil { panic(de) }
        fmt.Fprintf(w, "Writing Loader Datafile: %v\n", loaderHandler.Filename)
        fmt.Fprintf(w, "Writing Datafile: %v\n", dataHandler.Filename)
        io.Copy(lf, loaderFile)
        io.Copy(df, dataFile)
        defer lf.Close()
        defer df.Close()
      case  "pgrestore": dp.PGRestore(w)
      case  "psql":
        dataFile, dataHandler, dataErr := r.FormFile("data")
        if dataErr != nil { panic(dataErr) }
        dp.payloadFile = dataHandler.Filename
        df, de := os.OpenFile(os.Getenv("HOME")+"/psql/"+dataHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if de != nil { panic(de) }
        fmt.Fprintf(w, "Writing Datafile: %v\n", dataHandler.Filename)
        io.Copy(df, dataFile)
        defer df.Close()
      default: io.WriteString(w,"Utility is not matched\n")
    }
    io.WriteString(w, "Data File Uploaded\n")
    next.ServeHTTP(w, r)
  })
}

func (dp *Payload) DataLoader(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    switch dp.Utility {
      case "pgloader": dp.PGLoader(w)
      case  "pgrestore": dp.PGRestore(w)
      case  "psql": dp.PGsql(w)
      default: io.WriteString(w,"Utility is not matched\n")
    }
    next.ServeHTTP(w, r)
  })
}
