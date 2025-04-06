/server.go
package main

import (
  "fmt"
  "net/http"
   "strings"
)

const (
  localAddress string = "localhost:8082"
)

func echo(w http.ResponseWriter, req *http.Request) {
  fmt.Println("[ECHO] ...")
  for k, v := range req.Header {
    fmt.Printf("%s: %v\n", k, strings.Join(v, ", "))
  }
}

func main() {
 fmt.Println("Start mock server")

 http.HandleFunc("/", echo)
 http.ListenAndServe(localAddress, nil)
}
