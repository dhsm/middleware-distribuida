package main

import . "./message"
import . "./packet"
import . "./client_request_handler"

import "strconv"
import "fmt"


func main(){

  crh := ClientRequestHandler{}
  err := crh.NewCRH("tcp", "localhost", "8082", false, "JKDASBDK088a1asd")

  if (err != nil){
    return
  }

  i := 0 

  for{
    fmt.Println("sending packet",i)
    msg := Message{}
    msg.CreateMessage(strconv.Itoa(i), "replay", i, "client01")

    pkt := Packet{}
    if( i%2 == 1){
      // params := []string{}
      params := []string{"arg0", "arg1"}
      pkt.CreatePacket(ACK.Ordinal(), 0, params, msg)
    }else{
      params := []string{"arg0", "arg1", "arg2"}
      pkt.CreatePacket(MESSAGE.Ordinal(), 0, params, msg)  
    }
    

    crh.Send(pkt)
    response, err := crh.Receive()

    if(err != nil){
      panic(err)
    }

    // fmt.Println(response.GetMessage().MsgText)
    i, err = strconv.Atoi(response.GetMessage().MsgText)

    if(err != nil){
      panic(err)
    }
  }
}
