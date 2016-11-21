package main

import . "./message"
import . "./packet"
import . "./server_request_handler"
import "fmt"

func main(){
  srh := ServerRequestHandler{}
  srh.NewSRH("tcp", "127.0.0.1:8081")

  pktReceived := srh.Receive()
  fmt.Println(pktReceived)

  println("Creating response message...")
  msg := Message{}
  msg.CreateMessage("Hi Client", 99)

  println("Creating response packet...")
  pkt := Packet{}
  params := []string{}
  pkt.CreatePacket(MESSAGE, 0, params, msg)
  srh.Send(pkt)
}
