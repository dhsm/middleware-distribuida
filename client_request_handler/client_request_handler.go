package client_request_handler

import "net"
import "log"
import "sync"
import "io"
import "encoding/json"
import "encoding/gob"

import . "../packet"
import . "../message"

type OnPacketReceived func(pkt Packet)

type ClientRequestHandler struct {
	sync.Mutex
	Connection net.Conn
	Closed bool
}

func (crh *ClientRequestHandler) NewCRH(protocol string, host string, port string, isSubscriber bool, clientID string) error{
	gob.Register(Packet{})
	conn, err := net.Dial(protocol, net.JoinHostPort(host, port))
	crh.Connection = conn
	crh.Closed = false

	if (err != nil){
		log.Print("Error stablishing connection to: ", host, " at port ", port, " using ", protocol, " protocol.")
		return err
	}


	err = crh.SendType(isSubscriber, clientID)

	if (err != nil){
		log.Print("Error when registering a Client Request Handler for client: ", clientID)
		return err
	}

	crh.WaitAck(isSubscriber)

	if (err != nil){
		log.Print("Error when confirming registration of a Client Request Handler for client: ", clientID)
		return err
	}

	return nil
}

func (crh ClientRequestHandler) Send(pkt Packet) error{
	crh.Lock()
	encoded, err := json.Marshal(pkt)
	encoded_size, err := json.Marshal(len(encoded))
	crh.Connection.Write(encoded_size)
	crh.Connection.Write(encoded)
	crh.Unlock()

	if (err != nil){
		log.Print("Encoding error sending packet", err)
	}
	return err
}

func (crh ClientRequestHandler) Receive () (Packet, error){
	var pkt Packet
	var masPktSize int64

	crh.Lock()
	size := make([]byte, 3)
	io.ReadFull(crh.Connection,size)
	_ = json.Unmarshal(size, &masPktSize)
	packetMsh := make([]byte, masPktSize)
	io.ReadFull(crh.Connection,packetMsh)
	_ = json.Unmarshal(packetMsh, &pkt)
	crh.Unlock()

	return pkt, nil
}

func (crh ClientRequestHandler) SendAndReceive (pkt Packet) (Packet, error){
	err := crh.Send(pkt)
	if (err != nil){
		return Packet{}, nil
	}
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
		log.Print("Error receiving ACK: ", err)
	}

	if((isSubscriber && pkt.GetType() == REGISTER_RECEIVER_ACK) || (!isSubscriber && pkt.GetType() == REGISTER_SENDER_ACK)){
		return nil
	}

	log.Print("Ack for register not received")

	return err
}

func (crh ClientRequestHandler) Close() error{
	crh.Lock()
	err := crh.Connection.Close()
	crh.Unlock()
	if (err != nil){
		log.Print("Erro closing connection. ", err)
	}


	return err
}

func (crh ClientRequestHandler) SendAsync(pkt Packet){
	go func (){
		crh.Lock()
		encoded, err := json.Marshal(pkt)
		encoded_size, err := json.Marshal(len(encoded))
		crh.Connection.Write(encoded_size)
		crh.Connection.Write(encoded)
		crh.Unlock()

		if (err != nil){
			log.Print("Encoding error sending packet", err)
		}
	}()
}

func (crh ClientRequestHandler) ListenIncomingPackets(listener OnPacketReceived){
	// go func () {
	// 	for !crh.Closed{
	// 		pkt, err := crh.Receive()
	//
	// 		if (err!=nil){
	// 			crh.Closed = true
	// 		}else{
	// 			crh.Lock()
	// 			//TODO call connection
	// 			//crh.Connection.OnPacketReceived(pkt)
	// 			crh.Unlock()
	// 		}
	// 	}
	// }()
}

func (crh ClientRequestHandler) SetConnection(connection net.Conn){
	crh.Lock()
	crh.Connection = connection
	crh.Unlock()
}
