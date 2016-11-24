package broker

import "net"
import "errors"
import "time"
import "log"

import . "../client_request_handler"
import . "../message"

type HandleType uint

const (
	UNKNOWN Operation = iota
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
	Server Server
	ID int
	Connection net.Conn
	Running bool
	Type HandleType
	ClientID string
	ToSend chan Packet
	ClientID string
	WaitingACK WaitingACKSafe
}

func (ch *ConnectionHandler) NewCH(id int, conn net.Conn, server Server){
	ch.ID = id
	ch.Connection = conn
	ch.Running = false
	ch.Type = UNKNOWN
	ch.ToSend = make(chan Packet, 50)
	ch.WaitingACK{}
	ch.WaitingACK.Init()
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
	io.ReadFull(ch.Connection,size)
	_ = json.Unmarshal(size, &masPktSize)
	packetMsh := make([]byte, masPktSize)
	io.ReadFull(ch.Connection,packetMsh)
	_ = json.Unmarshal(packetMsh, &pkt)

	return pkt, nil
}

func (ch *ConnectionHandler) HandleRegister () error {
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
	pkt := Packet{}
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
	pkt := Packet{}
	params := []string{ch.ClientID}
	pkt.CreatePacket(REGISTER_RECEIVER_ACK, 0, params, Message{})
	err := ch.Send(pkt) 

	ch.Server.HandleRegisterReceiver(pkt, ch.ID)

	return err
}

func (ch *ConnectionHandler) Execute () {
	err := ch.HandleRegister()
	if(err != nil){
		ch.Running = false
	}

	for ch.Running{
		if(ch.Type == SENDER){
			err = ch.HandleReceivedMessages()
			if(err != nil){
				ch.Running = false
			}
		}else if(ch.Type == RECEIVER){
			err = ch.SendMessages()
			if(err != nil){
				ch.Running = false
			}
		}
	}

	ch.Connection.Close()

}

func (ch *ConnectionHandler) HandleReceivedMessages () error{
	pkt, err := ch.Receive()

	if(err != nil){
		return err
	}

	if(pkt.IsSubscribe()){
		ch.Server.HandleSubscribe(pkt)
		return nil
	}

	if(pkt.IsUnsubscribe()){
		ch.Server.HandleUnsubscribe(pkt)
		return nil
	}

	if(pkt.IsCreateTopic()){
		ch.Server.HandleCreateTopic(pkt)
		return nil
	}

	if(pkt.IsMessage()){
		ch.Server.HandleMessage(pkt)
		return nil
	}

	if(pkt.IsACK()){

	}
}

func (ch *ConnectionHandler) SendMessages () error{
	pkt := <-ch.ToSend
	err := ch.Send(pkt)

	if(pkt.IsMessage()){
		log.Print("Sending message ", pkt.GetMessageID()," to client " , ch.ClientID)
		msg = pkt.GetMessage()
		ch.Server.Receivers[ch.ClientID].GetWaitingAck().Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()), msg.MessageID})
		log.Print("Received ack [id: " , pkt.GetMessageID() , "] [size: " , ch.Server.Receivers[ch.ClientID].GetWaitingAck().Len(), "]")
	}
}

func (ch *ConnectionHandler) handleACK (pkt Packet){
	log.Print("Received ack [id: " , pkt.GetMessageID() , "] [size: " , ch.Server.Receivers[ch.ClientID].GetWaitingAck().Len(), "]")
	ch.Server.Receivers[ch.ClientID].GetWaitingAck().Remove(pkt.GetID())
	log.Print("Received ack [id: " , pkt.GetMessageID() , "] [size: " , ch.Server.Receivers[ch.ClientID].GetWaitingAck().Len(), "]")
}

func (ch *ConnectionHandler) GetWaitingAck () WaitingACK{
	return ch.WaitingACK
}
