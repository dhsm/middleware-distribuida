package main

import . "./client_request_handler"
import "fmt"

func main(){
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1:8081")

  //sample := []byte{0,0,0,0,0}

  //crh.Send(sample)

  msgReceived := crh.Receive()

  fmt.Println(msgReceived)
}
