package broker

import "net"
// import "io"
import "errors"
import "time"
import "log"
import "sync"
import "fmt"
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
	Lock sync.Mutex
	ACK sync.Mutex
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
	println("==> ConnectionHandle created!")

	ch.ID = id
	ch.Connection = conn
	ch.Server = server
	ch.Running = true
	ch.Type = UNKNOWN
	ch.ToSend = make(chan Packet, 50)
	ch.WaitingACK = WaitingACKSafe{}
	ch.WaitingACK.Init()

}

func (ch *ConnectionHandler) Execute() {
	println("@@@ ConnectionHandler [EXECUTE]")

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

func (ch *ConnectionHandler) HandleReceivedMessages() error{
	println("@@@ ConnectionHandler handle[RECEIVED_MESSAGES]")

	pkt, err := ch.Receive()

	//println("@@@ ConnectionHandler BACK TO handle[RECEIVED_MESSAGES]")
	if(err != nil){
		return err
	}

	if(pkt.IsSubscribe()){
		println("@@@ ConnectionHandler packet was::: SUBSCRIBE")
		ch.Server.HandleSubscribe(pkt)
		return nil
	}else if(pkt.IsUnsubscribe()){
		println("@@@ ConnectionHandler packet was::: UNSUBSCRIBE")
		ch.Server.HandleUnsubscribe(pkt)
		return nil
	}else if(pkt.IsCreateTopic()){
		println("@@@ ConnectionHandler packet was::: CREATE_TOPIC")
		ch.Server.HandleCreateTopic(pkt)
		return nil
	}else if(pkt.IsMessage()){
		println("@@@ ConnectionHandler packet was::: MESSAGE")
		ch.Server.HandleMessage(pkt)
		return nil
	}else if(pkt.IsACK()){
		println("@@@ ConnectionHandler packet was::: ACK")
		return nil
	}
	println("@@@ ConnectionHandler packet was::: !!!NO PACKET!!!")
	return nil
}

func (ch *ConnectionHandler) handleACK (pkt Packet){
	println("@@@ ConnectionHandler handle[ACK]")


	//log.Print("Received ack for message ", pkt.GetMessage().MessageID, " from client ", pkt.GetClientID())
	//temp := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	//log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , temp.Len(), "]")
	//temp_3 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	//temp_3.Remove(pkt.GetMessage().MessageID)
	//temp_2 := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
	//log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , temp_2.Len(), "]")
}

func (ch *ConnectionHandler) SendMessages() error{
	defer ch.ACK.Unlock()
	ch.ACK.Lock()
	println("@@@ ConnectionHandler send[MESSAGES]")
	pkt := <-ch.ToSend
	//println("094328409238409823094 new message received")
	err := ch.Send(pkt)

	if(pkt.IsMessage()){
		//log.Print("Sending message ", pkt.GetMessage().MessageID," to client " , pkt.GetClientID())
		msg := pkt.GetMessage()
		wack := ch.Server.Receivers[ch.ClientID].GetWaitingAck()
		wack.Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()) +(5*1000), msg.MessageID})
		//log.Print("Received ack [id: " , pkt.GetID() , "] [size: " , wack.Len(), "]")
	}
	return err
}

func (ch *ConnectionHandler) HandleRegister() error {
	println("@@@ ConnectionHandler handle[REGISTER]")

	pkt, err := ch.Receive()

	if(err != nil){
		log.Print("Error receiving Registration ", err)
		return err
	}

	if (pkt.IsRegisterSender()){
		ch.ClientID = pkt.GetClientID()
		return ch.HandleRegisterSender(pkt)
	}else if(pkt.IsRegisterReceiver()){
		ch.ClientID = pkt.GetClientID()
		return ch.HandleRegisterReceiver(pkt)
	}

	return errors.New("Expecting a REGISTER_CONSUMER or REGISTER_PRODUCER packet.")
}

func (ch *ConnectionHandler) HandleRegisterReceiver (pkt Packet) error {
	println("@@@ ConnectionHandler handle[REGISTER_RECEIVER]")
	//log.Print("Producer registered ", ch.ID)
	ch.Type = RECEIVER

	//Creating ACK response
	pkt = Packet{}
	params := []string{ch.ClientID}
	pkt.CreatePacket(REGISTER_RECEIVER_ACK.Ordinal(), 0, params, Message{})
	err := ch.Send(pkt)

	ch.Server.HandleRegisterReceiver(pkt, ch.ID)

	mm := MessageManager{}
	mm.NewMM(ch.Server, ch.ClientID, &ch.ACK)
	go mm.Execute()

	return err
}

func (ch *ConnectionHandler) HandleRegisterSender (pkt Packet) error {
	println("@@@ ConnectionHandler handle[REGISTER_SENDER]")

	//log.Print("Sender registered ", ch.ID)
	ch.Type = SENDER

	//Creating ACK response
	pkt = Packet{}
	params := []string{ch.ClientID}
	pkt.CreatePacket(REGISTER_SENDER_ACK.Ordinal(), 0, params, Message{})
	err := ch.Send(pkt)

	ch.Server.HandleRegisterSender(pkt, ch.ID)

	return err
}

func (ch *ConnectionHandler) Send(pkt Packet) error{
	println("@@@ ConnectionHandler [SEND]")

	defer ch.Lock.Unlock()
	ch.Lock.Lock()

	encoded, err := json.Marshal(pkt)

	// println("√√√√√√√√√√√√")
	// println(len(encoded))
	// println("√√√√√√√√√√√√")

	if (err != nil){
		log.Print("Encoding error sending packet", err)
		return err
	}

	pkt_len := fmt.Sprintf("%06d",len(encoded))

	encoded_size, err := json.Marshal(pkt_len)

	if (err != nil){
		log.Print("Encoding error sending packet", err)
		return err
	}

	ch.Connection.Write(encoded_size)
	ch.Connection.Write(encoded)

	if (err != nil){
		log.Print("Encoding error sending packet", err)
		return err
	}

	// println("##########")
	// fmt.Println(pkt)
	// println("##########")

	return nil
}

func (ch *ConnectionHandler) Receive() (Packet, error){
	println("@@@ ConnectionHandler [RECEIVE]")

	var pkt Packet
	var masPktSize string
	var size int64

	defer ch.Lock.Unlock()
	ch.Lock.Lock()

	// println("@@@ ConnectionHandler [RECEIVE] *we are inside Lock()*")
	temp := make([]byte, 8)
	// println("@@@ ConnectionHandler [RECEIVE] *about to ReadFull*")
	read_len, err := ch.Connection.Read(temp)

	if(read_len != 8){
		panic("Deu merda!")
	}

	if(err != nil){
		log.Print(err)
	}
	// println("@@@ ConnectionHandler [RECEIVE] *ReadFull, now will Unmarshall*")
	err = json.Unmarshal(temp[:read_len], &masPktSize)

	if(err!=nil){
		log.Print(err)
	}

	fmt.Sscanf(masPktSize, "%d", &size)

	// println("√√√√√√√√√√√√")
	// println(size)
	// println("√√√√√√√√√√√√")

	packetMsh := make([]byte, size)
	// println("@@@ ConnectionHandler [RECEIVE] *read size packet*")
	read_len, err = ch.Connection.Read(packetMsh)

	if(err!=nil){
		log.Print(err)
	}
	// println("@@@ ConnectionHandler [RECEIVE] *read message packet*")
	err = json.Unmarshal(packetMsh[:read_len], &pkt)

	if(err!=nil){
		log.Print(err)
	}
	// println("@@@ ConnectionHandler [RECEIVE] *we will return now*")
	// println("##########")
	// fmt.Println(pkt)
	// println("##########")
	return pkt, err
}

func (ch *ConnectionHandler) GetWaitingAck() WaitingACKSafe{
	return ch.WaitingACK
}
