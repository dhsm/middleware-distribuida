package main

import . "./message"
import . "./server_request_handler"
//import "encoding/json"
import "fmt"

func main(){
  srh := ServerRequestHandler{}
  srh.NewSRH("tcp", "127.0.0.1:8081")

  msgReceived := srh.Receive()
  fmt.Println(msgReceived)

  // _ = json.Unmarshal(bytes, &msgUnmarshaled)
  // fmt.Print("mensagem recebida: ")
  // fmt.Println(bytes)
  // fmt.Println(msgUnmarshaled.Msgtext)

  msg := Message{"oi cliente", 99, 0}
  //msgMarshaled, _ := json.Marshal(msg)
  srh.Send(msg)
}
