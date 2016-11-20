package server_request_handler

import "net"
//import "bufio"
// import "fmt"
// import "io"
// import "bytes"
import "encoding/gob"
import . "../message"
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

func (srh ServerRequestHandler) Send(msg Message) {
  enc := gob.NewEncoder(srh.Connection)
	err := enc.Encode(msg)
	if (err != nil){
		log.Fatal("Encoding error", err)
	}
}

func (srh *ServerRequestHandler) Receive() Message {
  dec := gob.NewDecoder(srh.Connection)
  var msg Message
  err := dec.Decode(&msg)
  if (err != nil){
		log.Fatal("Decoding error", err)
	}
  return msg
  // result :=
  // p := make([]byte, 25)
  //
  // buf := new(bytes.Buffer)
  // //var k int32
  // _, _ := io.ReadFull(srh.Connection,buf)
  // fmt.Println(buf)
  //
  // var size int32
  // buf := bytes.NewReader()
  // err = binary.Read(buf, binary.BigEndian, &size)
  // if err != nil {
  //     fmt.Println(err)
  //     return
  // }
  //
  // fmt.Println(size)
  //
  // p := make([]byte, 49)
  // _, _ = io.ReadFull(srh.Connection,p)
  // return p

}
