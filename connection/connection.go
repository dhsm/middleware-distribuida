package connection

import "sync"

import . "../packet"
// import . "../message"
import . "../client_request_handler"
// import . "../server_request_handler"
import . "../topic_session"

type onPacketReceived func(pkt Packet)

type SubscribedSafe struct{
	sync.Mutex
	Map map[string][]onPacketReceived
}

func (sd *SubscribedSafe) Init() {
	sd.Map = make(map[string][]onPacketReceived)
}

type WaitingACKSafe struct{
	sync.Mutex
	Map map[int]Packet
}

func (was *WaitingACKSafe) Init() {
	was.Map = make(map[int]Packet)
}

type Connection struct{
	sync.Mutex
	ClientId string
	HostIp string
	HostPort string
	HostProtocol string

	ReceiverConnection ClientRequestHandler
	SenderConnection ClientRequestHandler

	Subscribed SubscribedSafe
	WaitingACK WaitingACKSafe

	Sessions []TopicSession
}

func (cnn *Connection) CreateConnection(host_ip string, host_port string, host_protocol string){
	cnn.HostIp = host_ip
	cnn.HostPort = host_port
	cnn.HostProtocol = host_protocol
	cnn.WaitingACK.Init()
	cnn.Subscribed.Init()
	cnn.Sessions = make([]TopicSession, 50)
}