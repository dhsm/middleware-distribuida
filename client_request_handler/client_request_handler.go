package client_request_handler

import "net"
import "bufio"
import "fmt"

type ClientRequestHandler struct {
  Connection net.Conn
}

func (crh ClientRequestHandler) NewCRH(protocol string, address string){
  conn, _ := net.Dial(protocol, address)
  crh.Connection = conn
}

func (crh ClientRequestHandler) Send(msg []byte){
  fmt.Fprint(crh.Connection,msg,"\n")
}

func (crh ClientRequestHandler) Receive () []byte{
  bytes, _ := bufio.NewReader(crh.Connection).ReadBytes('\n')
  return bytes
}
