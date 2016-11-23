package server_request_handler

import "net"
import "log"
import "io"
import "encoding/json"

import . "../packet"

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
  encoded, err := json.Marshal(pkt)
  encoded_size, err := json.Marshal(len(encoded))
  srh.Connection.Write(encoded_size)
  srh.Connection.Write(encoded)
	if (err != nil){
		log.Fatal("Encoding error", err)
	}
}

func (srh *ServerRequestHandler) Receive() Packet {
  var pkt Packet
  var masPktSize int64

  size := make([]byte, 3)
  io.ReadFull(srh.Connection,size)
  _ = json.Unmarshal(size, &masPktSize)
  packetMsh := make([]byte, masPktSize)
  io.ReadFull(srh.Connection,packetMsh) 
  _ = json.Unmarshal(packetMsh, &pkt)

  return pkt
}
