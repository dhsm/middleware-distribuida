package packet

import "time"

import . "../message"

type Enum interface {
	Name() string
	Ordinal() int
	ValueOf() *[]string
}

type Operation uint

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
func (op Operation) Ordinal() uint {
	return uint(op)
}
func (op Operation) String() string {
	return operations[op]
}

func (op Operation) Values() *[]string {
	return &operations
}

type Packet struct{
	Operation Operation
	ID uint
	Params []string
	Msg Message
	TimeStamp int32
	Index int //This is necessary because we are using a PriorityQueue
}

func (pkt *Packet) CreatePacket(op Operation, id uint, params []string, msg Message){
	pkt.Operation = op
	pkt.ID = id
	pkt.Params = params
	pkt.Msg = msg
	pkt.Index = -1
	pkt.TimeStamp = int32(time.Now().Unix())
}

func (pkt *Packet) GetOperation() Operation{
	return pkt.Operation
}

func (pkt *Packet) SetOperation(operation Operation){
	pkt.Operation = operation
}

func (pkt *Packet) GetID() uint{
	return pkt.ID
}

func (pkt *Packet) SetID(id uint){
	pkt.ID = id
}

func (pkt *Packet) GetType() Operation{
	return pkt.Operation
}

func (pkt *Packet) GetMessage() Message{
	return pkt.Msg
}

func (pkt *Packet) SetMessage(msg Message){
	pkt.Msg = msg
}

func (pkt *Packet) GetParams() []string{
	return pkt.Params
}

func (pkt *Packet) SetParams(params []string){
	pkt.Params = params
}