package broker

import "net"
import "io"
import "errors"
import "time"
import "log"
import "sync"
import "encoding/json"

import . "../message"
import . "../packet"

type HandleType uint

const (
	UNKNOWN HandleType = iota
	SENDER
	RECEIVER)

var handletypes = []string{
	"UNKNOWN",
	"SENDER",
	"RECEIVER"}

func (ht HandleType) Name() string {
	return handletypes[ht]
}
func (ht HandleType) Ordinal() uint {
	return uint(ht)
}
func (ht HandleType) String() string {
	return handletypes[ht]
}

func (ht HandleType) Values() *[]string {
	return &handletypes
}

type ConnectionHandler struct{
	sync.Mutex
	ID int
	Server Server
	Connection net.Conn
	Running bool
	Type HandleType
	ToSend chan Packet
	WaitingACK WaitingACKSafe
	ClientID string
}

func (ch *ConnectionHandler) NewCH(id int, conn net.Conn, server Server){
	ch.ID = id
	ch.Connection = conn
	ch.Server = server
	ch.Running = true
	ch.Type = UNKNOWN
	ch.ToSend = make(chan Packet, 50)
	ch.WaitingACK = WaitingACKSafe{}
	ch.WaitingACK.Init()
	
}

func (ch *ConnectionHandler) Execute () {
	err := ch.HandleRegister()
	if(err != nil){
		ch.Running = false
	}

	for ch.Running{
		if(ch.Type == SENDER){
			err = ch.HandleReceivedMessages()
		}else if(ch.Type == RECEIVER){
			err = ch.SendMessages()
		}

		if(err != nil){
			log.Print(err)
			ch.Running = false
		}
	}

	ch.Connection.Close()
}

func (ch *ConnectionHandler) HandleReceivedMessages () error{
	println("*** ConnectionHandler HandleReceivedMessages")

	pkt, err := ch.Receive()

	if(err != nil){
		return err
	}

	if(pkt.IsSubscribe()){
		ch.Server.HandleSubscribe(pkt)
		return nil
	}else if(pkt.IsUnsubscribe()){
		ch.Server.HandleUnsubscribe(pkt)
		return nil
	}else if(pkt.IsCreateTopic()){
		ch.Server.HandleCreateTopic(pkt)
		return nil
	}else if(pkt.IsMessage()){
		ch.Server.HandleMessage(pkt)
		return nil
	}else if(pkt.IsACK()){
		return nil
	}
	return nil
}

func (ch *ConnectionHandler) handleACK (pkt Packet){
	println("*** ConnectionHandler handleACK")

	log.Print("Received ack for message ", pkt.GetMessage().MessageID, " from client ", pkt.GetClientID())
	temp := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , temp.Len(), "]")
	temp_3 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	temp_3.Remove(pkt.GetMessage().MessageID)
	temp_2 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , temp_2.Len(), "]")
}

func (ch ConnectionHandler) SendMessages () error{
	pkt := <-ch.ToSend
	err := ch.Send(pkt)

	if(pkt.IsMessage()){
		log.Print("Sending message ", pkt.GetMessage().MessageID," to client " , pkt.GetClientID())
		msg := pkt.GetMessage()
		wack := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
		wack.Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()), msg.MessageID})
		log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , wack.Len(), "]")
	}
	return err
}


func (ch *ConnectionHandler) Send(pkt Packet) error{

	defer ch.Unlock()
	ch.Lock()
	encoded, err := json.Marshal(pkt)
	encoded_size, err := json.Marshal(len(encoded))
	ch.Connection.Write(encoded_size)
	ch.Connection.Write(encoded)

	if (err != nil){
		log.Print("Encoding error sending packet", err)
	}
	return err
}

func (ch *ConnectionHandler) Receive () (Packet, error){
	var pkt Packet
	var masPktSize int64

	defer ch.Unlock()
	ch.Lock()
	size := make([]byte, 3)
	_, err := io.ReadFull(ch.Connection,size)
	if(err != nil){
		log.Print(err)
	}
	err = json.Unmarshal(size, &masPktSize)
	packetMsh := make([]byte, masPktSize)
	_, err = io.ReadFull(ch.Connection,packetMsh)
	if(err != nil){
		log.Print(err)
	}
	err = json.Unmarshal(packetMsh, &pkt)

	return pkt, err
}

func (ch *ConnectionHandler) HandleRegister () error {
	println("*** ConnectionHandler HandleRegister")
	
	pkt, err := ch.Receive()

	if(err != nil){
		log.Print("Error receiving Registration ", err)
		return err
	}

	if (pkt.IsRegisterSender()){
		ch.ClientID = pkt.GetClientID()
		return ch.HandleRegisterSender(pkt)
	}

	if(pkt.IsRegisterReceiver()){
		ch.ClientID = pkt.GetClientID()
		return ch.HandleRegisterReceiver(pkt)
	}

	return errors.New("Expecting a REGISTER_CONSUMER or REGISTER_PRODUCER packet.")
}

func (ch *ConnectionHandler) HandleRegisterSender (pkt Packet) error {
	log.Print("Sender registered ", ch.ID)
	ch.Type = SENDER

	//Creating ACK response
	pkt = Packet{}
	params := []string{ch.ClientID}
	pkt.CreatePacket(REGISTER_SENDER_ACK, 0, params, Message{})
	err := ch.Send(pkt)

	ch.Server.HandleRegisterSender(pkt, ch.ID)

	return err
}

func (ch *ConnectionHandler) HandleRegisterReceiver (pkt Packet) error {
	log.Print("Producer registered ", ch.ID)
	ch.Type = RECEIVER

	//Creating ACK response
	pkt = Packet{}
	params := []string{ch.ClientID}
	pkt.CreatePacket(REGISTER_RECEIVER_ACK, 0, params, Message{})
	err := ch.Send(pkt)

	ch.Server.HandleRegisterReceiver(pkt, ch.ID)

	return err
}

func (ch *ConnectionHandler) GetWaitingAck() WaitingACKSafe{
	return ch.WaitingACK
}
