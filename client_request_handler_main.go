package main

import . "./message"
import "encoding/json"
import . "./client_request_handler"
import "fmt"

func main(){
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1:8081", false)

  msg := Message{"oi servidor", 99, 0}

  msgMarshaled, _ := json.Marshal(msg)
  fmt.Print("Total na mensagem a ser enviada")
  fmt.Println(len(msgMarshaled))
  fmt.Print("Mensagem marshaled")
  fmt.Println(msgMarshaled)
  fmt.Println(msg.Msgtext)

  crh.Send(msg)

  var msgReceived Message
  msgReceived, err := crh.Receive()
  if(err != nil){
    fmt.Println(err)
  }
  fmt.Print("bytes da mensagem recebida: ")
  fmt.Println(msgReceived)
  // var msgUnmarshaled Message
  // _ = json.Unmarshal(msgReceived, &msgUnmarshaled)
  // fmt.Print("mensagem recebida: ")
  // fmt.Println(msgUnmarshaled)
}
