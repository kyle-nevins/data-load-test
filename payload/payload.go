package payload

import (
  "io"
  "fmt"
  "net/http"
  "os"
  "github.com/gorilla/mux"
  "encoding/json"
  "strconv"
)

type Payload struct {
  payloadFile   string
  Utility       string
  uri           string
  db_host       string
  db_port       string
  username      string
  db_name       string
  password      string
}

// Middleware function, which will be called for each request
func (dp *Payload) AssignVariables(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dp.Utility = vars["utility"]
    fmt.Fprintf(w, "Upload Utility %v\n", dp.Utility)

    //Set up the VCAP services and get the interface credentials.
    vcap := os.Getenv("VCAP_SERVICES")
  	var vsMap map[string][]map[string]interface{}
  	err := json.Unmarshal([]byte(vcap), &vsMap)
    if err != nil {	panic(err) }

    //Set up the environment credentials as a map.
    var credsInterface interface{}
    if _, exist := vsMap["postgresql-9.5-odb"]; exist {
      credsInterface = vsMap["postgresql-9.5-odb"][0]["credentials"]
    } else {
      credsInterface = vsMap["postgresql-9.5"][0]["credentials"]
    }

    //Pass the rest of the arguments to a variable to assign into the payload
    dp.uri = credsInterface.(map[string]interface{})["uri"].(string)
    dp.password = credsInterface.(map[string]interface{})["password"].(string)
    dp.db_host = credsInterface.(map[string]interface{})["db_host"].(string)
    dp.db_port = strconv.FormatFloat(credsInterface.(map[string]interface{})["db_port"].(float64), 'f', -1, 64)
    dp.username = credsInterface.(map[string]interface{})["username"].(string)
    dp.db_name = credsInterface.(map[string]interface{})["db_name"].(string)
}

func (dp *Payload) WriteFile(w http.ResponseWriter, r *http.Request) {
    switch dp.Utility {
      case "pgloader":
        loaderFile, loaderHandler, _ := r.FormFile("loaderFile")
        dataFile, dataHandler, _ := r.FormFile("data")
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
        dataFile, dataHandler, _ := r.FormFile("data")
        dp.payloadFile = dataHandler.Filename
        df, de := os.OpenFile(os.Getenv("HOME")+"/pg_restore/"+dataHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if de != nil { panic(de) }
        fmt.Fprintf(w, "Writing Datafile: %v\n", dataHandler.Filename)
        io.Copy(df, dataFile)
        defer df.Close()
      case  "psql":
        dataFile, dataHandler, _ := r.FormFile("data")
        dp.payloadFile = dataHandler.Filename
        df, de := os.OpenFile(os.Getenv("HOME")+"/psql/"+dataHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if de != nil { panic(de) }
        fmt.Fprintf(w, "Writing Datafile: %v\n", dataHandler.Filename)
        io.Copy(df, dataFile)
        defer df.Close()
      default: io.WriteString(w,"Utility is not matched\n")
    }
    io.WriteString(w, "Data File Uploaded\n")
}

func (dp *Payload) DataLoader(w http.ResponseWriter, r *http.Request) {
    switch dp.Utility {
      case "pgloader": dp.PGLoader(w)
      case  "pgrestore": dp.PGRestore(w)
      case  "psql": dp.PGsql(w)
      default: io.WriteString(w,"Utility is not matched\n")
  }
}
