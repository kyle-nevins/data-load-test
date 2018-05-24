package payload

import (
  "io"
  "os"
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"
  "os/exec"
  "bytes"
  "log"
)

func (dp Payload) PGLoader(w http.ResponseWriter) {
  io.WriteString(w,"Writing data using the PGLoader utility.\n")

  path := os.Getenv("HOME")+"/pgloader/loaderFile"
  read, err := ioutil.ReadFile(path)
	if err != nil { log.Panic(err) }

	newContents := strings.Replace(string(read),
    "postgresql://username:password@host:port/database",
    dp.uri,
    -1)

	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil { log.Panic(err) }

  command:= exec.Command(os.Getenv("HOME")+"/pgloader/pgloader",path)
  command.Dir = os.Getenv("HOME")+"/pgloader"
  var out bytes.Buffer
  command.Stdout = &out
  err = command.Run()
  if err != nil { log.Panic(err) }
  log.Println(out.String())
}

func (dp Payload) PGRestore(w http.ResponseWriter) {
  io.WriteString(w,"Writing data using the PGRestore utility.\n")

  //Command run for PGSQL
  command := exec.Command(os.Getenv("HOME")+"/pg_restore/pg_restore")

  env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", dp.password))
	command.Env = env

  //Pass the rest of the arguments to a variable to assign into the command.
  var args []string
  args = append(args,"--host", dp.db_host)
  args = append(args,"--port", dp.db_port)
  args = append(args,"--username", dp.username)
  args = append(args,"--dbname", dp.db_name)
  args = append(args,"--clean --jobs=10 --no-owner --no-privileges --verbose")
  args = append(args,os.Getenv("HOME")+"/pg_restore/"+dp.payloadFile)
  command.Args = args

  fmt.Fprintf(w,"%v %v",os.Getenv("HOME")+"/pg_restore/pg_restore",args)

  command.Dir = os.Getenv("HOME")+"/pg_restore"
  var out bytes.Buffer
  command.Stdout = &out
  err := command.Run()
  if err != nil { log.Println(err) }
  log.Println(out.String())
}

func (dp Payload) PGsql(w http.ResponseWriter) {

  //Intro line to write to the browser/command line.
  io.WriteString(w,"Writing data using the PGSQL utility.\n")

  //Command run for PGSQL
  command := exec.Command(os.Getenv("HOME")+"/psql/psql")

  env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", dp.password))
	command.Env = env

  fmt.Fprintf(w,"%v",command.Env)

  //Pass the rest of the arguments to a variable to assign into the command.
  var args []string
  args = append(args,"--host", dp.db_host)
  args = append(args,"--port", dp.db_port)
  args = append(args,"--username", dp.username)
  args = append(args,"--dbname", dp.db_name)
  args = append(args,"--file", os.Getenv("HOME")+"/psql/"+dp.payloadFile)
  command.Args = args

  fmt.Fprintf(w,"%v %v",os.Getenv("HOME")+"/psql/psql",args)

  command.Dir = os.Getenv("HOME")+"/psql"
  var out bytes.Buffer
  command.Stdout = &out
  err := command.Run()
  if err != nil { fmt.Println(err) }

  fmt.Fprintf(w,"%v",out.String())

}
