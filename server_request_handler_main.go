package main

import . "./message"
import . "./packet"
import . "./client_request_handler"

// import "time"
import "strconv"

func main(){
  srh := ClientRequestHandler{}
  srh.NewCRHS("tcp", ":8082")

  for{
    pktRCV, err := srh.Receive()

    if(err != nil){
      panic(err)
    }

    print(pktRCV.GetMessage().MsgText)
    ret, err := strconv.Atoi(pktRCV.GetMessage().MsgText)

    if(err != nil){
      panic(err)
    }

    ret++

    msg := Message{}
    msg.CreateMessage(strconv.Itoa(ret), "replay", ret, "server01")

    pkt := Packet{}
    params := []string{}
    pkt.CreatePacket(ACK.Ordinal(), ret-1, params, msg)

    srh.Send(pkt)
    // time.Sleep(time.Second*5)
  }
}
