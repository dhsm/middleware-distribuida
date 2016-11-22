package main

import . "./message"
import . "./packet"
import . "./client_request_handler"
import "fmt"

func main(){
  println("Starting Client Request Handler...")
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1", "8081", false, "JKDASBDK088a1asd")

  println("Creating message...")
  msg := Message{}
  msg.CreateMessage("Hi Server!", 99)
  
  println("Creating packet...")
  pkt := Packet{}
  params := []string{"arg0", "arg1", "arg2"}
  pkt.CreatePacket(MESSAGE, 0, params, msg)

  crh.Send(pkt)

  var pktReceived Packet
  pktReceived, err := crh.Receive()
  if(err != nil){
    fmt.Println(err)
  }
  fmt.Print("Packet received:")
  fmt.Println(pktReceived)
}
