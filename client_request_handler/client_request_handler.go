package client_request_handler

import "net"
import "log"
import "sync"
// import "io"
import "encoding/json"
import "fmt"
// import "math/rand"
// import "time"

import . "../packet"
import . "../message"

type ClientRequestHandler struct {
	S sync.Mutex
	R sync.Mutex
	Connection net.Conn
	Listen net.Listener
	Closed bool
	CNN PacketListener
	RemoteAddr string
	Protocol string
	Subscriber bool
	ClientID string
}

func (crh *ClientRequestHandler) NewCRH(protocol string, host string, port string, isSubscriber bool, clientID string) error{
	conn, err := net.Dial(protocol, net.JoinHostPort(host, port))

	crh.Connection = conn
	crh.Closed = false

	crh.RemoteAddr = net.JoinHostPort(host, port)
	crh.Protocol = protocol
	crh.Subscriber = isSubscriber


	if (err != nil){
		log.Fatal("Error stablishing connection to: ", host, " at port ", port, " using ", protocol, " protocol.")
		return err
	}

	// println("Sending client type...")
	err = crh.SendType(isSubscriber, clientID)

	if (err != nil){
		log.Fatal("Error when registering a Client Request Handler for client: ", clientID)
		return err
	}

	// println("Waiting ACK confirming...")
	crh.WaitACK(isSubscriber)

	if (err != nil){
		log.Fatal("Error when confirming registration of a Client Request Handler for client: ", clientID)
		return err
	}

	println("Setup completed!")

	return nil
}

func (crh *ClientRequestHandler) NewCRHS(protocol string, port string) error{
	ln, err := net.Listen(protocol, port)

	if (err != nil){
		log.Fatal(err)
		return err
	}

	conn, err := ln.Accept()

	if (err != nil){
		log.Fatal(err)
		return err
	}

	crh.Listen = ln
	crh.Connection = conn
	err = crh.ReceiveType()

	if(err!=nil){
		log.Fatal(err)
		return err
	}

	err = crh.SendACK()

	if(err!=nil){
		log.Fatal(err)
		return err
	}

	return err

}

func (crh *ClientRequestHandler) Send(pkt Packet) error{
	println("ΩΩΩ ClientRequestHandler send[PACKET]")
	defer crh.S.Unlock()
	crh.S.Lock()

	encoded, err := json.Marshal(pkt)

	if (err != nil){
		log.Fatal(err)
		return err
	}

	// println("√√√√√√√√√√√√")
	// println(len(encoded))
	// println("√√√√√√√√√√√√")

	pkt_len := fmt.Sprintf("%06d",len(encoded))

	println(pkt_len)

	encoded_size, err := json.Marshal(pkt_len)

	if (err != nil){
		log.Fatal(err)
		return err
	}

	crh.Connection.Write(encoded_size)
	crh.Connection.Write(encoded)

	// println("Packet sent sucessfully!")
	// println("##########")
	// fmt.Println(pkt)
	// println("##########")

	return nil
}

func (crh *ClientRequestHandler) Receive() (Packet, error){
	println("ΩΩΩ ClientRequestHandler receive[PACKET]")
	defer crh.R.Unlock()
	crh.R.Lock()

	var pkt Packet
	var masPktSize string
	var size int64

	temp := make([]byte, 8)

	read_len, err := crh.Connection.Read(temp)

	if(err!=nil){
		log.Fatal(err)
	}

	err = json.Unmarshal(temp[:read_len], &masPktSize)

	if(err!=nil){
		log.Fatal(err)
	}

	fmt.Sscanf(masPktSize, "%d", &size)

	// println("√√√√√√√√√√√√")
	// println(size)
	// println("√√√√√√√√√√√√")

	packetMsh := make([]byte, size)

	read_len, err = crh.Connection.Read(packetMsh)

	if(err!=nil){
		log.Fatal(err)
	}

	err = json.Unmarshal(packetMsh[:read_len], &pkt)

	if(err!=nil){
		log.Fatal(err)
	}

	// println("Packet received sucessfully!")
	// println("##########")
	// fmt.Println(pkt)
	// println("##########")

	return pkt, err
}

func (crh *ClientRequestHandler) SendAndReceive(pkt Packet) (Packet, error){
	err := crh.Send(pkt)
	if (err != nil){
		return Packet{}, nil
	}
	return crh.Receive()
}

func (crh *ClientRequestHandler) SendType(isSubscriber bool, clientID string) error{
	pkt := Packet{}
	msg := Message{}
	params := []string{clientID}
	if (isSubscriber){
		pkt.CreatePacket(REGISTER_RECEIVER.Ordinal(), 0, params, msg)
	}else{
		pkt.CreatePacket(REGISTER_SENDER.Ordinal(), 0, params, msg)
	}
	return crh.Send(pkt)
}

func (crh *ClientRequestHandler) ReceiveType() error{
	pkt, err := crh.Receive()

	if(pkt.IsRegisterReceiver()){
		crh.Subscriber = true
	}else{
		crh.Subscriber = false
	}

	crh.ClientID = pkt.GetParams()[0]
	return err
}

func (crh *ClientRequestHandler) SendACK() error{
	pkt := Packet{}
	msg := Message{}
	params := []string{}
	if (crh.Subscriber){
		pkt.CreatePacket(REGISTER_RECEIVER_ACK.Ordinal(), 0, params, msg)
	}else{
		pkt.CreatePacket(REGISTER_SENDER_ACK.Ordinal(), 0, params, msg)
	}
	return crh.Send(pkt)
}

func (crh *ClientRequestHandler) WaitACK(isSubscriber bool) error{
	pkt, err := crh.Receive()

	if (err != nil){
		log.Fatal("Error receiving ACK: ", err)
	}

	if((isSubscriber && pkt.GetType() == REGISTER_RECEIVER_ACK.Ordinal()) || (!isSubscriber && pkt.GetType() == REGISTER_SENDER_ACK.Ordinal())){
		return nil
	}

	// log.Fatal("Ack for register not received")

	return err
}

func (crh *ClientRequestHandler) Close() error{
	defer crh.S.Unlock()
	defer crh.R.Unlock()
	crh.S.Lock()
	crh.R.Lock()
	if(crh.Connection != nil){
		err := crh.Connection.Close()
		if (err != nil){
			log.Fatal("Erro closing connection. ", err)
		}
		return err
	}

	if(crh.Connection != nil){
		err := crh.Connection.Close()
		if (err != nil){
			log.Fatal("Erro closing connection. ", err)
		}
		return err
	}

	return nil
}

func (crh *ClientRequestHandler) SendAsync(pkt Packet){
	go func(){
		// println("ΩΩΩ ClientRequestHandler sendAsync[PACKET]")
		crh.S.Lock()

		encoded, err := json.Marshal(pkt)

		if (err != nil){
			log.Fatal(err)
		}

		// println("√√√√√√√√√√√√")
		// println(len(encoded))
		// println("√√√√√√√√√√√√")

		pkt_len := fmt.Sprintf("%06d",len(encoded))

		// println(pkt_len)

		encoded_size, err := json.Marshal(pkt_len)

		if (err != nil){
			log.Fatal(err)
		}

		crh.Connection.Write(encoded_size)
		crh.Connection.Write(encoded)

		// println("Packet sent sucessfully!")
		// println("##########")
		// fmt.Println(pkt)
		// println("##########")
		crh.S.Unlock()
	}()
}

func (crh *ClientRequestHandler) ListenIncomingPackets(){
	go func() {
		for !crh.Closed{
			pkt, err := crh.Receive()

			if (err!=nil){
				crh.Closed = true
			}else{
				crh.CNN.OnPacket(pkt)
			}
		}
	}()
}

func (crh *ClientRequestHandler) SetConnection(connection PacketListener){
	crh.R.Lock()
	crh.CNN = connection
	crh.R.Unlock()
}

func (crh *ClientRequestHandler) GetClientID() string{
	return crh.ClientID
}

func (crh *ClientRequestHandler) IsSubscriber() bool{
	return crh.Subscriber
}
