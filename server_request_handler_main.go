package main

import . "./message"
import . "./packet"
import . "./server_request_handler"
import "fmt"
// import "time"

func main(){
  srh := ServerRequestHandler{}
  srh.NewSRH("tcp", "127.0.0.1:8081")

  pktReceived := srh.Receive()
  fmt.Println(pktReceived)

  
  msg := Message{}

  println("Creating response packet...")
  pkt := Packet{}
  params := []string{}
  if(pktReceived.GetType() == REGISTER_SENDER){
    pkt.CreatePacket(REGISTER_SENDER_ACK, 0, params, msg)  
  }else if(pktReceived.GetType() == REGISTER_RECEIVER){
    pkt.CreatePacket(REGISTER_RECEIVER_ACK, 0, params, msg)  
  }else{
    pkt.CreatePacket(MESSAGE, 0, params, msg)
  }
  srh.Send(pkt)

  pktReceived = srh.Receive()
  println("Creating response message...")
  msg.CreateMessage("Hi Client", "notopic", 99, "semid")
  pkt.CreatePacket(MESSAGE, 0, params, msg)
  srh.Send(pkt)

  println("Vou dormir")
  // time.Sleep(time.Second * 5)

  println("Acordei!")

  for i := 0;; i++ {
    println("Sending packet...")
    pkt.CreatePacket(MESSAGE, uint(i), params, msg)
    srh.Send(pkt)
  }
}
