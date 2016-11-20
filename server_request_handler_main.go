package main

import . "./message"
import . "./server_request_handler"
import "encoding/json"
import "fmt"

func main(){
  srh := ServerRequestHandler{}
  srh.NewSRH("tcp", "127.0.0.1:8081")

  bytes := srh.Receive()

  var msgUnmarshaled Message
  _ = json.Unmarshal(bytes, &msgUnmarshaled)
  fmt.Print("mensagem recebida: ")
  fmt.Println(msgUnmarshaled.Msgtext)

  msg := Message{"oi cliente", 99, 0}
  msgMarshaled, _ := json.Marshal(msg)
  srh.Send(msgMarshaled)
}
