package client_request_handler

import "net"
//import "bufio"
import "errors"
import . "../message"
import "encoding/gob"
import "log"

type ClientRequestHandler struct {
	Connection net.Conn
	Async bool
}

func (crh *ClientRequestHandler) NewCRH(protocol string, address string, async bool){
	//conn, err := net.Dial(protocol, address)
	conn, _ := net.Dial(protocol, address)
	crh.Connection = conn
	crh.Async = async
}

func (crh ClientRequestHandler) Send(msg Message){
	enc := gob.NewEncoder(crh.Connection)
	err := enc.Encode(msg)
	if (err != nil){
		log.Fatal("Encoding error", err)
	}
	// j := int32(len(msg))
	// buf := new(bytes.Buffer)
	// err := binary.Write(buf, binary.BigEndian, j)
	// if err != nil {
	// 		fmt.Println(err)
	// 		return
	// }
	//
	// binary.Write(buf,binary.BigEndian, msg)
	//
	// crh.Connection.Write(buf)
	// if (crh.Async){ //This is a async call and we don't need to keep the conection alive waiting for an answer
	// 	crh.Connection.Close()
	// }
}

func (crh ClientRequestHandler) Receive () (Message, error){
	if (crh.Async){ //If the call was async we are not expecting a response right away so a error will be issued
		var msg Message
		return msg, errors.New("Connection already closed!")
	}
	// bytes, _ := bufio.NewReader(crh.Connection).ReadBytes('\n')
	// crh.Connection.Close()
	// return bytes, nil
	dec := gob.NewDecoder(crh.Connection)
  var msg Message
  err := dec.Decode(&msg)
  if (err != nil){
		log.Fatal("Decoding error", err)
	}
  return msg, nil
}
