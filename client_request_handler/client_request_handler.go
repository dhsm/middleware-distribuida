package client_request_handler

import "net"
import "log"
import "sync"
import "io"
// import "fmt"
import "encoding/json"
// import "math/rand"
// import "time"
import "strconv"
import "strings"

import . "../packet"
import . "../message"

type ClientRequestHandler struct {
	S sync.Mutex
	R sync.Mutex
	ConnectionSend net.Conn
	ConnectionReceive net.Conn
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

	// if (err != nil){
	// 	log.Fatal(err)
	// 	return err
	// }

	crh.ConnectionSend = conn
	crh.ConnectionReceive = conn
	crh.Closed = false

	// r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// tries := 50
	// rport := 0
	// for tries > 0 {
	// 	rport = int((r.Int()%10000)+2000)
	// 	// println("Listerning on...")
	// 	// println(strings.Split(conn.LocalAddr().String(),":")[0]+":"+strconv.Itoa(rport))
	// 	ln, err := net.Listen(protocol, strings.Split(conn.LocalAddr().String(),":")[0]+":"+strconv.Itoa(rport))

	// 	tries--

	// 	if(err == nil && ln != nil){
	// 		tries = 0
	// 		crh.Listen = ln
	// 		break
	// 	}else{
	// 		log.Fatal(err)
	// 	}
	// }
	
	crh.RemoteAddr = net.JoinHostPort(host, port)
	crh.Protocol = protocol

	// encoded_port, err := json.Marshal(rport)
	// crh.ConnectionSend.Write(encoded_port)
	// crh.ConnectionSend.Close()

	if (err != nil){
		log.Fatal("Error stablishing connection to: ", host, " at port ", port, " using ", protocol, " protocol.")
		return err
	}

	println("Sending client type...")
	err = crh.SendType(isSubscriber, clientID)

	if (err != nil){
		log.Fatal("Error when registering a Client Request Handler for client: ", clientID)
		return err
	}

	println("Waiting ACK confirming...")
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
	crh.ConnectionReceive = conn
	crh.ConnectionSend = conn
	crh.Protocol = protocol

	rport := 0

	size := make([]byte, 4)
	
	_, err = io.ReadFull(crh.ConnectionReceive,size)
	
	if(err!=nil){
		log.Fatal(err)
		return err
	}
	
	err = json.Unmarshal(size, &rport)
	
	if(err!=nil){
		log.Fatal(err)
		return err
	}

	crh.RemoteAddr = strings.Split(conn.RemoteAddr().String(),":")[0]+":"+strconv.Itoa(rport)
	// crh.ConnectionReceive.Close()
	
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

func (crh ClientRequestHandler) Send(pkt Packet) error{
	crh.S.Lock()

	println("Sengind new packet to", crh.RemoteAddr)

	// conn, err := net.Dial(crh.Protocol, crh.RemoteAddr)

	// println(crh.Protocol)

	// if (err != nil){
	// 	log.Fatal(err)
	// 	return err
	// }

	// crh.ConnectionSend = conn
	
	encoded, err := json.Marshal(pkt)
	
	if (err != nil){
		log.Fatal(err)
		return err
	}

	encoded_size, err := json.Marshal(len(encoded))
	
	if (err != nil){
		log.Fatal(err)
		return err
	}

	crh.ConnectionSend.Write(encoded_size)
	crh.ConnectionSend.Write(encoded)
	// crh.ConnectionSend.Close()

	crh.S.Unlock()

	println("Packet sent sucessfully!")

	return nil
}

func (crh ClientRequestHandler) Receive() (Packet, error){
	crh.R.Lock()

	// conn, err := crh.Listen.Accept()

	// if(err!=nil){
	// 	log.Fatal(err)
	// }

	// crh.ConnectionReceive = conn

	println("Esperando pacote de", crh.ConnectionReceive.RemoteAddr().String())


	var pkt Packet
	var masPktSize int64

	
	size := make([]byte, 3)
	
	_, err := io.ReadFull(crh.ConnectionReceive,size)
	
	if(err!=nil){
		log.Fatal(err)
	}
	
	err = json.Unmarshal(size, &masPktSize)
	
	if(err!=nil){
		log.Fatal(err)
	}
	
	packetMsh := make([]byte, masPktSize)
	
	_, err = io.ReadFull(crh.ConnectionReceive,packetMsh)

	// crh.ConnectionReceive.Close()
	
	crh.R.Unlock()

	if(err!=nil){
		log.Fatal(err)
	}

	err = json.Unmarshal(packetMsh, &pkt)

	if(err!=nil){
		log.Fatal(err)
	}	

	return pkt, nil
}

func (crh ClientRequestHandler) SendAndReceive(pkt Packet) (Packet, error){
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

func (crh ClientRequestHandler) ReceiveType() error{
	pkt, err := crh.Receive()

	if(pkt.IsRegisterReceiver()){
		crh.Subscriber = true
	}else{
		crh.Subscriber = false
	}

	crh.ClientID = pkt.GetParams()[0]
	return err
}

func (crh ClientRequestHandler) SendACK() error{
	pkt := Packet{}
	msg := Message{}
	params := []string{}
	if (crh.Subscriber){
		pkt.CreatePacket(REGISTER_RECEIVER_ACK, 0, params, msg)
	}else{
		pkt.CreatePacket(REGISTER_SENDER_ACK, 0, params, msg)
	}
	return crh.Send(pkt)
}

func (crh ClientRequestHandler) WaitACK(isSubscriber bool) error{
	pkt, err := crh.Receive()

	if (err != nil){
		log.Fatal("Error receiving ACK: ", err)
	}

	if((isSubscriber && pkt.GetType() == REGISTER_RECEIVER_ACK) || (!isSubscriber && pkt.GetType() == REGISTER_SENDER_ACK)){
		return nil
	}

	log.Fatal("Ack for register not received")

	return err
}

func (crh ClientRequestHandler) Close() error{
	defer crh.S.Unlock()
	defer crh.R.Unlock()
	crh.S.Lock()
	crh.R.Lock()
	if(crh.ConnectionSend != nil){
		err := crh.ConnectionSend.Close()
		if (err != nil){
			log.Fatal("Erro closing connection. ", err)
		}
		return err
	}

	if(crh.ConnectionReceive != nil){
		err := crh.ConnectionReceive.Close()
		if (err != nil){
			log.Fatal("Erro closing connection. ", err)
		}
		return err
	}

	return nil
}

func (crh ClientRequestHandler) SendAsync(pkt Packet){
	go func(){
		crh.S.Lock()
		encoded, err := json.Marshal(pkt)
		encoded_size, err := json.Marshal(len(encoded))
		crh.ConnectionSend.Write(encoded_size)
		crh.ConnectionSend.Write(encoded)
		crh.S.Unlock()

		if (err != nil){
			log.Fatal("Encoding error sending packet", err)
		}
	}()
}

func (crh ClientRequestHandler) ListenIncomingPackets(){
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

func (crh ClientRequestHandler) SetConnection(connection PacketListener){
	crh.R.Lock()
	crh.CNN = connection
	crh.R.Unlock()
}

func (crh ClientRequestHandler) GetClientID() string{
	return crh.ClientID
}

func (crh ClientRequestHandler) IsSubscriber() bool{
	return crh.Subscriber
}