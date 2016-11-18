package main

import . "./message"
import "encoding/json"
import . "./client_request_handler"
import "fmt"

func main(){
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1:8081")

  msg := Message{"oi servidor"}

  msgMarshaled, _ := json.Marshal(msg)
  fmt.Print("mensagem em bytes: ")
  fmt.Println(len(msgMarshaled))
  crh.Send(msgMarshaled)

  msgReceived := crh.Receive()
  fmt.Print("bytes da mensagem recebida: ")
  fmt.Println(msgReceived)
  var msgUnmarshaled Message
  _ = json.Unmarshal(msgReceived, &msgUnmarshaled)
  fmt.Print("mensagem recebida: ")
  fmt.Println(msgUnmarshaled)
}
