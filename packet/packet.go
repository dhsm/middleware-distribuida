package packet

import . "../message"

type Enum interface {
	name() string
	ordinal() int
	valueOf() *[]string
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

func (op Operation) name() string {
	return operations[op]
}
func (op Operation) ordinal() uint {
	return uint(op)
}
func (op Operation) String() string {
	return operations[op]
}

func (op Operation) values() *[]string {
	return &operations
}

type PacketHeader struct{
	Operation Operation
	Id uint
}

func (ph *PacketHeader) getOperation() Operation{
	return ph.Operation
}

func (ph *PacketHeader) setOperation(operation Operation){
	ph.Operation = operation
}

func (ph *PacketHeader) getId() uint{
	return ph.Id
}

func (ph *PacketHeader) setId(id uint){
	ph.Id = id
}

type PacketBody struct{
	Params []string
	Msg Message
}

func (pb *PacketBody) getMessage() Message{
	return pb.Msg
}

func (pb *PacketBody) setMessage(msg Message){
	pb.Msg = msg
}

func (pb *PacketBody) getParams() []string{
	return pb.Params
}

func (pb *PacketBody) setParams(params []string){
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

func (pkt *Packet) getBody() PacketBody{
	return pkt.Body
}

func (pkt *Packet) setBody(body PacketBody){
	pkt.Body = body
}

func (pkt *Packet) getHeader() PacketHeader{
	return pkt.Header
}

func (pkt *Packet) setHeader(header PacketHeader){
	pkt.Header = header
}