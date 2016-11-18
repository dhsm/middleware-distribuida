package client_request_handler

import "net"
import "bufio"

type ClientRequestHandler struct {
  Connection net.Conn
}

func (crh *ClientRequestHandler) NewCRH(protocol string, address string){
  //conn, err := net.Dial(protocol, address)
  conn, _ := net.Dial("tcp", "127.0.0.1:8081")
  crh.Connection = conn
}

func (crh ClientRequestHandler) Send(msg []byte){
  crh.Connection.Write(msg)
}

func (crh ClientRequestHandler) Receive () []byte{
  bytes, _ := bufio.NewReader(crh.Connection).ReadBytes('\n')
  //var buf bytes
  crh.Connection.Close()
  return bytes
}
