package main

import . "./message"
import "encoding/json"
import . "./client_request_handler"
import "fmt"

func main(){
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1:8081")

  msg := Message{"oi servidor"}

  msgMarshaled, err := json.Marshal(msg)
  fmt.Println(err)
  fmt.Println(len(msgMarshaled))
  fmt.Println(msgMarshaled)
  //sample := []byte{0,0,0,0,1}
  //crh.Send(sample)
  crh.Send(msgMarshaled)

  msgReceived := crh.Receive()
  fmt.Println(msgReceived)
  var msgUnmarshaled Message
  _ = json.Unmarshal(msgReceived, &msgUnmarshaled)
  fmt.Println(err)
  fmt.Println(msgUnmarshaled)
}
