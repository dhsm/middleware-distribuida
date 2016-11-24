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
	Server Server
	ID int
	Connection net.Conn
	Running bool
	Type HandleType
	ClientID string
	ToSend chan Packet
	WaitingACK WaitingACKSafe
}

func (ch *ConnectionHandler) NewCH(id int, conn net.Conn, server Server){
	ch.ID = id
	ch.Connection = conn
	ch.Running = false
	ch.Type = UNKNOWN
	ch.ToSend = make(chan Packet, 50)
	ch.WaitingACK = WaitingACKSafe{}
	ch.WaitingACK.Init()
	ch.Server = server
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
	_, err = io.ReadFull(ch.Connection,size)
	if(err != nil){
		return err
	}
	err := json.Unmarshal(size, &masPktSize)
	packetMsh := make([]byte, masPktSize)
	_, err = io.ReadFull(ch.Connection,packetMsh)
	if(err != nil){
		return err
	}
	err = json.Unmarshal(packetMsh, &pkt)

	return pkt, err
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
		return nil
	}
	return nil
}

func (ch ConnectionHandler) SendMessages () error{
	pkt := <-ch.ToSend
	err := ch.Send(pkt)

	if(pkt.IsMessage()){
		log.Print("Sending message ", pkt.GetMessage()," to client " , ch.ClientID)
		msg := pkt.GetMessage()
		//ch.Server.Receivers[ch.ClientID].GetWaitingAck().Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()), msg.MessageID})
		teste := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
		teste.Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()), msg.MessageID})
		//log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , ch.Server.Receivers[ch.ClientID].GetWaitingAck().Len(), "]")
		teste_2 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
		log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , teste_2.Len(), "]")
	}
	return err
}

func (ch *ConnectionHandler) handleACK (pkt Packet){
	temp := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , temp.Len(), "]")
	temp_3 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	temp_3.Remove(pkt.GetMessage().MessageID)
	temp_2 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , temp_2.Len(), "]")
}

func (ch *ConnectionHandler) GetWaitingAck() WaitingACKSafe{
	return ch.WaitingACK
}
