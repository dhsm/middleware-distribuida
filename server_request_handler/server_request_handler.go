package server_request_handler

import "net"
//import "bufio"
// import "fmt"
// import "io"
// import "bytes"
import "encoding/gob"
import . "../packet"
import "log"

type ServerRequestHandler struct{
  Listen net.Listener
  Connection net.Conn
}

func (srh *ServerRequestHandler) NewSRH(protocol string, port string){
  ln, _ := net.Listen(protocol, port)
  conn, _ := ln.Accept()

  srh.Listen = ln
  srh.Connection = conn
}

func (srh ServerRequestHandler) Send(pkt Packet) {
  enc := gob.NewEncoder(srh.Connection)
	err := enc.Encode(pkt)
	if (err != nil){
		log.Fatal("Encoding error", err)
	}
}

func (srh *ServerRequestHandler) Receive() Packet {
  dec := gob.NewDecoder(srh.Connection)
  var pkt Packet
  err := dec.Decode(&pkt)
  if (err != nil){
		log.Fatal("Decoding error", err)
	}
  return pkt
}
