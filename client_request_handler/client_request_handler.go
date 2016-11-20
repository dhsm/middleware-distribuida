package client_request_handler

import "net"
import "bufio"
import "errors"

type ClientRequestHandler struct {
	Connection net.Conn
	Async bool
}

func (crh *ClientRequestHandler) NewCRH(protocol string, address string, async bool){
	//conn, err := net.Dial(protocol, address)
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	crh.Connection = conn
	crh.Async = async
}

func (crh ClientRequestHandler) Send(msg []byte){
	crh.Connection.Write(msg)
	if (crh.Async){ //This is a async call and we don't need to keep the conection alive waiting for an answer
		crh.Connection.Close()
	}
}

func (crh ClientRequestHandler) Receive () ([]byte, error){
	if (crh.Async){ //If the call was async we are not expecting a response right away so a error will be issued
		var bytes []byte
		return bytes, errors.New("Connection already closed!")
	}
	bytes, _ := bufio.NewReader(crh.Connection).ReadBytes('\n')
	crh.Connection.Close()
	return bytes, nil
}
