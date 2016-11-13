package main

import . "./message"
import . "./client_request_handler"
import "encoding/json"
import "fmt"

func main(){
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1:8081")

  msg := Message{"oi servidor"}
  msgMarshaled, _ := json.Marshal(msg)

  crh.Send(msgMarshaled)

  msgReceived := crh.Receive()

  var msgUnmarshaled Message
  message := json.Unmarshal(msgReceived, &msgUnmarshaled)
  fmt.Println(message)
}
