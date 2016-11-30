package packet

import "time"
import "fmt"

import . "../message"

type Operation int

const (
	REGISTER_SENDER Operation = iota
	REGISTER_RECEIVER
	REGISTER_SENDER_ACK
	REGISTER_RECEIVER_ACK
	SUBSCRIBE
	UNSUBSCRIBE
	CREATE_TOPIC
	MESSAGE
	ACK)

var operations = []string{
	"REGISTER_SENDER",
	"REGISTER_RECEIVER",
	"REGISTER_SENDER_ACK",
	"REGISTER_RECEIVER_ACK",
	"SUBSCRIBE",
	"UNSUBSCRIBE",
	"CREATE_TOPIC",
	"MESSAGE",
	"ACK"}

func (op Operation) Name() string {
	return operations[op]
}
func (op Operation) Ordinal() int {
	return int(op)
}
func (op Operation) String() string {
	return operations[op]
}

func (op Operation) Values() *[]string {
	return &operations
}

type Packet struct{
	Operation int
	ID int
	Params []string
	TimeStamp int32
	Index int //This is necessary because we are using a PriorityQueue

	//Message content
	MsgText string
	Priority int
	Destination string
	MessageID string
	Redelivery bool
}

func (pkt Packet) String() string {
	return fmt.Sprintf("Operation: %d\nID: %d\nParams: %p\nTimeStamp: %d\nIndex: %d\nMsgText: %s\nPriority: %d\nDestination: %s\nMessageID: %s\nRedelivery: %t\n", pkt.Operation, pkt.ID, pkt.Params, pkt.TimeStamp, pkt.Index, pkt.MsgText, pkt.Priority, pkt.Destination, pkt.MessageID, pkt.Redelivery)
}

func (pkt *Packet) CreatePacket(op int, id int, params []string, msg Message){
	pkt.Operation = op
	pkt.ID = id
	pkt.Params = params
	pkt.Index = -1
	pkt.TimeStamp = int32(time.Now().Unix())

	//Incorporating message into packet
	pkt.MsgText = msg.MsgText
	pkt.Priority = msg.Priority
	pkt.Destination = msg.Destination
	pkt.MessageID = msg.MessageID
	pkt.Redelivery = msg.Redelivery
}

func (pkt *Packet) GetOperation() int{
	return pkt.Operation
}

func (pkt *Packet) SetOperation(operation int){
	pkt.Operation = operation
}

func (pkt *Packet) GetID() int{
	return pkt.ID
}

func (pkt *Packet) GetClientID() string{
	return pkt.Params[0]
}

func (pkt *Packet) SetID(id int){
	pkt.ID = id
}

func (pkt *Packet) GetType() int{
	return pkt.Operation
}

func (pkt *Packet) GetMessage() Message{
	return Message{pkt.MsgText, pkt.Priority, pkt.Destination, pkt.MessageID, pkt.Redelivery}
}

func (pkt *Packet) SetMessage(msg Message){
	pkt.MsgText = msg.MsgText
	pkt.Priority = msg.Priority
	pkt.Destination = msg.Destination
	pkt.MessageID = msg.MessageID
	pkt.Redelivery = msg.Redelivery
}

func (pkt *Packet) GetParams() []string{
	return pkt.Params
}

func (pkt *Packet) SetParams(params []string){
	pkt.Params = params
}

func (pkt *Packet) IsRegisterSender() bool{
	return pkt.Operation == REGISTER_SENDER.Ordinal()
}
func (pkt *Packet) IsRegisterReceiver() bool{
	return pkt.Operation == REGISTER_RECEIVER.Ordinal()
}
func (pkt *Packet) IsRegisterSenderAck() bool{
	return pkt.Operation == REGISTER_SENDER_ACK.Ordinal()
}
func (pkt *Packet) IsRegisterReceiverAck() bool{
	return pkt.Operation == REGISTER_RECEIVER_ACK.Ordinal()
}
func (pkt *Packet) IsSubscribe() bool{
	return pkt.Operation == SUBSCRIBE.Ordinal()
}
func (pkt *Packet) IsUnsubscribe() bool{
	return pkt.Operation == UNSUBSCRIBE.Ordinal()
}
func (pkt *Packet) IsCreateTopic() bool{
	return pkt.Operation == CREATE_TOPIC.Ordinal()
}
func (pkt *Packet) IsMessage() bool{
	return pkt.Operation == MESSAGE.Ordinal()
}
func (pkt *Packet) IsACK() bool{
	return pkt.Operation == ACK.Ordinal()
}