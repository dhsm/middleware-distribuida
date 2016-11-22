package client_request_handler

import . "../packet"
import . "../message"
import "net"
import "encoding/gob"
import "log"

type ClientRequestHandler struct {
	Connection net.Conn
}

func (crh *ClientRequestHandler) NewCRH(protocol string, host string, port string, isSubscriber bool, clientID string) error{
	conn, err := net.Dial(protocol, net.JoinHostPort(host, port))

	crh.Connection = conn

	if (err != nil){
		log.Fatal("Error stablishing connection to: ", host, " at port ", port, " using ", protocol, " protocol.")
		return err
	}


	err = crh.SendType(isSubscriber, clientID)
	
	if (err != nil){
		log.Fatal("Error when registering a Client Request Handler for client: ", clientID)
		return err
	}

	crh.WaitAck(isSubscriber)

	if (err != nil){
		log.Fatal("Error when confirming registration of a Client Request Handler for client: ", clientID)
		return err
	}

	return nil
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

func (crh ClientRequestHandler) SendAndReceive (pkt Packet) (Packet, error){
	crh.Send(pkt)
	return crh.Receive()
}

func (crh ClientRequestHandler) SendType(isSubscriber bool, clientID string) error{
	pkt := Packet{}
	msg := Message{}
	params := []string{clientID}
	if (isSubscriber){
		pkt.CreatePacket(REGISTER_RECEIVER, 0, params, msg)	
	}else{
		pkt.CreatePacket(REGISTER_SENDER, 0, params, msg)	
	}
	return crh.Send(pkt)
}

func (crh ClientRequestHandler) WaitAck(isSubscriber bool) error{
	pkt, err := crh.Receive()

	if (err != nil){
		log.Fatal("Error receiving ACK", err)
	}

	if((isSubscriber && pkt.GetType() == REGISTER_RECEIVER_ACK) || (!isSubscriber && pkt.GetType() == REGISTER_SENDER_ACK)){
		return nil
	}

	log.Fatal("Ack for register not received")

	return err
}