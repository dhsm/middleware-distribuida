package connection

import "../packet"
import "../message"
import "../client_request_handler"
import "../server_request_handler"
import "net"

type Connection struct{
	CHR ClientRequestHandler
	WaitingACK map[int]Packet
	BrokerAddress string
	BrokerPort string
	BrokerProtocol string
	InputStream <-chan Message
	OutputStream chan<- Message
}

func (cnn *Connection) CreateConnection(br_addr string, br_port string, br_protocol string){
	cnn.BrokerAddress = br_addr
	cnn.BrokerPort = br_port
	cnn.BrokerProtocol = br_protocol
	cnn.InputStream = make(chan Packet, 250)
	cnn.OutputStream = make(chan Packet, 250)
	cnn.WaitingACK = make(map[int]Packet)
}

func (cnn *Connection) GetInputStream() <-chan Message{
	return cnn.InputStream
}

func (cnn *Connection) GetOutputStream() chan<- Message{
	return cnn.OutputStream
}
>>>>>>> a36972906546bd3ed359de998c89dbc6a86bff60
