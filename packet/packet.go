package packet

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

type PacketHeader struct{
	Operation Operation
	Id uint
}

func (ph *PacketHeader) GetOperation() Operation{
	return ph.Operation
}

func (ph *PacketHeader) SetOperation(operation Operation){
	ph.Operation = operation
}

func (ph *PacketHeader) GetId() uint{
	return ph.Id
}

func (ph *PacketHeader) SetId(id uint){
	ph.Id = id
}

type PacketBody struct{
	Params []string
	Msg Message
}

func (pb *PacketBody) GetMessage() Message{
	return pb.Msg
}

func (pb *PacketBody) SetMessage(msg Message){
	pb.Msg = msg
}

func (pb *PacketBody) GetParams() []string{
	return pb.Params
}

func (pb *PacketBody) SetParams(params []string){
	pb.Params = params
}

type Packet struct{
  Header PacketHeader
  Body PacketBody
}

func (pkt *Packet) CreatePacket(op Operation, id uint, params []string, msg Message){
  pkt.Header = PacketHeader{op, id}
  pkt.Body = PacketBody{params, msg}
}

func (pkt *Packet) GetBody() PacketBody{
	return pkt.Body
}

func (pkt *Packet) SetBody(body PacketBody){
	pkt.Body = body
}

func (pkt *Packet) GetHeader() PacketHeader{
	return pkt.Header
}

func (pkt *Packet) SetHeader(header PacketHeader){
	pkt.Header = header
}