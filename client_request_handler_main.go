package main

import . "./message"
import . "./packet"
import . "./middleware"
import . "./client_request_handler"
import "fmt"
import "time"

func printReceivedPackets(pkt Packet){
  fmt.Println(pkt)
}

func main(){
  println("Starting Client Request Handler...")
  crh := ClientRequestHandler{}
  err := crh.NewCRH("tcp", "127.0.0.1", "8081", false, "JKDASBDK088a1asd")

  if (err != nil){
    return
  }

  println("Creating message...")
  msg := Message{}
  msg.CreateMessage("Hi Server!", "notopic", 99, "semid")
  
  println("Creating packet...")
  pkt := Packet{}
  params := []string{"arg0", "arg1", "arg2"}
  pkt.CreatePacket(MESSAGE, 0, params, msg)

  crh.SendAsync(pkt)
  // time.Sleep(time.Second)
  // crh.SendAsync(pkt)

  var pktReceived Packet
  pktReceived, err = crh.Receive()
  if(err != nil){
    fmt.Println(err)
  }
  fmt.Print("Packet received:")
  fmt.Println(pktReceived)
  crh.SendAsync(pkt)
  crh.SendAsync(pkt)
  crh.SendAsync(pkt)
  cnn := Connection{}
  cnn.CreateConnection("127.0.0.1", "8081", "tcp")
  crh.SetConnection(cnn)
  crh.ListenIncomingPackets()
  println("Waiting 20 before end execution...")
  time.Sleep(time.Second * 20)
}
