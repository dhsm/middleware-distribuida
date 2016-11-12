package server_request_handler

import "net"
import "bufio"

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

func (srh ServerRequestHandler) Send(msg []byte) {
  srh.Connection.Write(msg)
}

func (srh ServerRequestHandler) Receive() []byte {
  bytes, _ := bufio.NewReader(srh.Connection).ReadBytes('\n')
  return bytes
}
