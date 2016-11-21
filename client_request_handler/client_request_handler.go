package client_request_handler

import . "../packet"
import "net"
import "encoding/gob"
import "log"

type ClientRequestHandler struct {
	Connection net.Conn
}

func (crh *ClientRequestHandler) NewCRH(protocol string, host string, port string, bool isSubscriber, String clientID){
	crh.Connection, _ = net.Dial(protocol, net.JoinHostPort(host, port))
}

func (crh ClientRequestHandler) Send(pkt Packet) error{
	enc := gob.NewEncoder(crh.Connection)
	err := enc.Encode(pkt)
	if (err != nil){
		log.Fatal("Encoding error sending packet", err)
	}
	return err
}

func (crh ClientRequestHandler) Receive () (Packet, error){
	var pkt Packet
	dec := gob.NewDecoder(crh.Connection)
	err := dec.Decode(&pkt)
	if (err != nil){
		log.Fatal("Decoding error", err)
	}
	return pkt, err
}

func (crh ClientRequestHandler) SendType(bool isSubscriber, string clientID) {
	pkt := Packet{}
	msg := Message{}
	params := []string{clientID}
	if (isSubscriber){
		pkt.CreatePacket(REGISTER_RECEIVER, 0, params, msg)	
	}else{
		pkt.CreatePacket(REGISTER_SENDER, 0, params, msg)	
	}
	crh.send(pkt)
}