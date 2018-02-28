package payload

import (
  "io"
  "os"
  "fmt"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "strings"
  "os/exec"
  "bytes"
  "strconv"
)

func (dp Payload) PGLoader(w http.ResponseWriter) {
  io.WriteString(w,"Writing data using the PGLoader utility.\n")

  path := os.Getenv("HOME")+"/pgloader/loaderFile"
  read, err := ioutil.ReadFile(path)
	if err != nil { panic(err) }

	newContents := strings.Replace(string(read),
    "postgresql://username:password@host:port/database",
    os.Getenv("DATABASE_URL"),
    -1)

	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {	panic(err) }

  command:= exec.Command(os.Getenv("HOME")+"/pgloader/pgloader",path)
  command.Dir = os.Getenv("HOME")+"/pgloader"
  var out bytes.Buffer
  command.Stdout = &out
  err = command.Run()
  if err != nil { fmt.Println(err) }

  fmt.Fprintf(w,"%v",out.String())
}

func (dp Payload) PGRestore(w http.ResponseWriter) {
  io.WriteString(w,"Writing data using the PGRestore utility.\n")
}

func (dp Payload) PGsql(w http.ResponseWriter) {

  //Intro line to write to the browser/command line.
  io.WriteString(w,"Writing data using the PGSQL utility.\n")

  //Set up the VCAP services and get the interface credentials.
  vcap := os.Getenv("VCAP_SERVICES")
	var vsMap map[string][]map[string]interface{}
	err := json.Unmarshal([]byte(vcap), &vsMap)
  if err != nil {	panic(err) }

  //Command run for PGSQL
  command := exec.Command(os.Getenv("HOME")+"/psql/psql")

  //Set up the password as an environment variable.
  credsInterface := vsMap["postgresql-9.5-odb"][0]["credentials"]
  password := credsInterface.(map[string]interface{})["password"].(string)
  env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", password))
	command.Env = env

  //Pass the rest of the arguments to a variable to assign into the command.
  var args []string
  args = append(args,"--host ", credsInterface.(map[string]interface{})["db_host"].(string))
  args = append(args,"--port ", strconv.FormatFloat(credsInterface.(map[string]interface{})["db_port"].(float64), 'f', -1, 64))
  args = append(args,"--username ", credsInterface.(map[string]interface{})["username"].(string))
  args = append(args,"--dbname ", credsInterface.(map[string]interface{})["db_name"].(string))
  args = append(args,"--f ", os.Getenv("HOME")+"/psql/"+dp.payloadFile)
  command.Args = args

  fmt.Fprintf(w,"%v %v",os.Getenv("HOME")+"/psql/psql",args)

  command.Dir = os.Getenv("HOME")+"/psql"
  var out bytes.Buffer
  command.Stdout = &out
  err = command.Run()
  if err != nil { fmt.Println(err) }

  fmt.Fprintf(w,"%v",out.String())

}
