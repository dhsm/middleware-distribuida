package main

import . "./message"
import . "./server_request_handler"
import "encoding/json"
import "fmt"

func main(){
  srh := ServerRequestHandler{}
  srh.NewSRH("tcp", ":8081")

  bytes := srh.Receive()

  var msgUnmarshaled Message
  message := json.Unmarshal(bytes, &msgUnmarshaled)
  fmt.Println(message)
}
