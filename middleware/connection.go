package middleware

import "sync"
import "math"
import "time"
import "log"
import "errors"
import "reflection"

import "github.com/nu7hatch/gouuid"

import . "../packet"
import . "../message"
import . "../client_request_handler"

type OnMessageReceived func(msg Message)

type SubscribedSafe struct{
	sync.Mutex
	Map map[string][]OnMessageReceived
}

func (sd *SubscribedSafe) Init() {
	sd.Map = make(map[string][]OnMessageReceived)
}

func (sd *SubscribedSafe) Get(key string) ([]OnMessageReceived, bool){
	defer sd.Unlock()
	sd.Lock()
	l, found := sd.Map[key]
	return l, found
}

func (sd *SubscribedSafe) Add(key string, fu OnMessageReceived){
	defer sd.Unlock()
	sd.Lock()
	_, f := sd.Map[key]
	if (!f) {
		sd.Map[key] = make([]OnMessageReceived,0)
	}

	//Checking if this listener is not already in this list
	for _, fun := range sd.Map[key] {
		if(reflect.DeepEqual(fun, fu)){
			return
		}
	}

	sd.Map[key] = append(sd.Map[key], fu)
}

func (sd *SubscribedSafe) Remove(key string, fu OnMessageReceived){
	f, e := sd.Map[key]

	if(!e){
		return false
	}

	//Removing the function from the list of listeners for the topic
	for i, elem := range f {
		if(reflect.DeepEqual(elem, fu)){
			f[i] = f[len(f)-1]
			f[len(f)-1] = nil
			f = f[:len(f)-1]
			return true
		}
	}

	return false
}

type MessageWaitingAck struct{
	Message Message
	TimeStamp int32
	MessageID string
}

type WaitingACKSafe struct{
	sync.Mutex
	Map map[string]MessageWaitingAck
}

func (was *WaitingACKSafe) Init() {
	was.Map = make(map[string]MessageWaitingAck)
}

func (was *WaitingACKSafe) Get(key string) (MessageWaitingAck, bool){
	defer was.Unlock()
	was.Lock()
	l, found := was.Map[key]
	return l, found
}

func (was *WaitingACKSafe) Add(key string, msg MessageWaitingAck){
	defer was.Unlock()
	was.Lock()
	was.Map[key] = msg
}

func (was *WaitingACKSafe) Remove(key string){
	defer was.Unlock()
	was.Lock()
	delete(was.Map,key)
}

func (was *WaitingACKSafe) Len(key string) int{
	defer was.Unlock()
	was.Lock()
	return len(was.Map)
}

func (was *WaitingACKSafe) Peek(key string) (string, MessageWaitingAck, bool){
	defer was.Unlock()
	was.Lock()
	for k, e := range was.Map {
		return k, e, true
	}

	return nil, MessageWaitingAck{}, false
}


type Connection struct{
	Lock sync.Mutex
	MessageSent sync.Cond
	AckReceived sync.Cond

	ClientID string
	HostIp string
	HostPort string
	HostProtocol string

	ReceiverConnection ClientRequestHandler
	SenderConnection ClientRequestHandler

	Subscribed SubscribedSafe
	WaitingACK WaitingACKSafe

	Sessions []TopicSession

	Stopped bool
	Open bool
	Modified bool

	PacketIDGenerator uint
}

func (cnn *Connection) CreateConnection(host_ip string, host_port string, host_protocol string){
	cnn.Lock = sync.Mutex{}
	cnn.MessageSent = sync.Cond{L: &cnn.Lock}
	cnn.AckReceived = sync.Cond{L: &cnn.Lock}

	cnn.HostIp = host_ip
	cnn.HostPort = host_port
	cnn.HostProtocol = host_protocol
	cnn.WaitingACK.Init()
	cnn.Subscribed.Init()
	cnn.Sessions = make([]TopicSession, 0)
	uuid_, _ := uuid.NewV4()
	cnn.ClientID = uuid_.String()
	cnn.Stopped = true
	cnn.Open = false
	cnn.Modified = false
	cnn.PacketIDGenerator = 0
}

func (cnn *Connection) IsOpen() error{
	if (!cnn.Open){
		return errors.New("Operation not allowed in closed connection.")
	}
	return nil
}

func (cnn *Connection) SetModified(){
	cnn.Modified = true
}

func (cnn *Connection) GetClientID() string{
	return cnn.ClientID
}

func (cnn *Connection) SetClientID(clientid string) error{
	if(cnn.Modified){
		return errors.New("Change the client id is not allowed afthe the connection has been modified.")
	}

	cnn.ClientID = clientid
	return nil
}

func (cnn *Connection) Close() error{

	return nil
}

func (cnn *Connection) CreateSession() TopicSession{
	cnn.SetModified()
	tp := TopicSession{}
	//TODO: Add create session call when ready
	cnn.Sessions = append(cnn.Sessions, tp)
	return tp
}

func (cnn *Connection) SendMessage(msg Message) error{
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}

	cnn.WaitingACK.Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()), msg.MessageID})

	cnn.Lock.Lock()
	//Broadcasting that there is new messages waiting for an ACK
	cnn.MessageSent.Broadcast()
	cnn.Lock.Unlock()

	cnn.SetModified()

	pkt := Packet{}
	pkt.CreatePacket(MESSAGE, cnn.PacketIDGenerator, nil, msg)
	cnn.PacketIDGenerator++
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) SubscribeSessionToDestination(topic Topic, fu OnMessageReceived){
	defer cnn.Lock.Unlock()
	cnn.Lock.Lock()
	cnn.Subscribed.Add(topic.GetTopicName(), fu)
}

func (cnn *Connection) UnsubscribeSessionToDestination(topic Topic, fu OnMessageReceived) bool{
	defer cnn.Lock.Unlock()
	cnn.Lock.Lock()
	cnn.Subscribe.Remove(topic.GetTopicName(), fu)
}

func (cnn *Connection) Subscribe(topic Topic, fu OnMessageReceived) error{
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	cnn.SetModified()
	cnn.SubscribeSessionToDestination(topic, fu)
	pkt := Packet{}
	params := []string{cnn.ClientID, topic.GetTopicName()}
	pkt.CreatePacket(SUBSCRIBE, cnn.PacketIDGenerator, params, Message{})
	cnn.PacketIDGenerator++
	return cnn.SenderConnection.Send(pkt)
}

func (cnn *Connection) Unsubscribe(topic Topic, fu OnMessageReceived) error{
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	cnn.SetModified()
	result := cnn.UnsubscribeSessionToDestination(topic, fu)
	if (result){
		pkt := Packet{}
		params := []string{cnn.ClientID, topic.GetTopicName()}
		pkt.CreatePacket(UNSUBSCRIBE, cnn.PacketIDGenerator, params, Message{})
		cnn.PacketIDGenerator++
		return cnn.SenderConnection.Send(pkt)
	}

	return nil
}

func (cnn *Connection) AcknowledgeMessage(msg Message, ts TopicSession) error{
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	cnn.SetModified()
	pkt := Packet{}
	params := []string{cnn.ClientID, msg.MessageID}
	pkt.CreatePacket(ACK, cnn.PacketIDGenerator, params, Message{})
	cnn.PacketIDGenerator++
	return cnn.SenderConnection.SendAsync(pkt)
}

func (cnn *Connection) CloseSession(ts TopicSession){
	// err := cnn.IsOpen()
	// if(err != nil){
	// 	log.Print(err)
	// 	return err
	// }
	// cnn.SetModified()
	// pkt := Packet{}
	// params := []string{cnn.ClientID, msg.MessageID}
	// pkt.CreatePacket(ACK, cnn.PacketIDGenerator, params, Message{})
	// cnn.PacketIDGenerator++
	// return cnn.SenderConnection.SendAsync(pkt)
}

func (cnn *Connection) CreateTopic(tp Topic) error{
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	pkt := Packet{}
	params := []string{cnn.ClientID, tp.GetTopicName()}
	pkt.CreatePacket(CREATE_TOPIC, cnn.PacketIDGenerator, params, Message{})
	cnn.PacketIDGenerator++
	return cnn.SenderConnection.SendAsync(pkt)
}

func (cnn *Connection) ProcessACKS(){
	go func () {
		for{
			err := cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}

			if (cnn.WaitingACK.Len() == 0){
				cnn.Lock.Lock()

				err := cnn.IsOpen()
				if(err != nil){
					log.Print(err)
					break
				}

				cnn.MessageSent.Wait() //Waiting for messages to be sent before stat to process ACKS again
				cnn.Lock.Unlock()
			}

			err := cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}

			key, maa, f := cnn.WaitingACK.Peek()
			if(f){
				curr := int32(time.Now().Unix())
				if(maa.TimeStamp <= curr){
					cnn.WaitingACK.Remove(key)
					cnn.SendMessage(maa.Message)
				}else{
					time.Sleep(time.Microsecond * Time.Duration(curr))
				}
			}

			err := cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}
		}
	}()
}

func (cnn *Connection) OnPacketReceived(pkt Packet){

}

func (cnn *Connection) Start(){
	if(!cnn.Open){
		tries := 0
		for tries <= 5 {
			tries++
			errr := cnn.ReceiverConnection.NewCRH(cnn.HostProtocol, cnn.HostIp, cnn.HostPort, false, cnn.GetClientID())
			errs := cnn.SenderConnection.NewCRH(cnn.HostProtocol, cnn.HostIp, cnn.HostPort, false, cnn.GetClientID())

			if(errr != nil || errs != nil){
				cnn.ReceiverConnection.Close()
				cnn.SenderConnection.Close()
				delay := math.Pow(2,float64(tries))
				time.Sleep(time.Second * time.Duration(delay))
				log.Print("Error stablishing connection for client ", cnn.ClientID, ", trying again in ", delay, " seconds...")
				continue
			}

			cnn.Open = true
			cnn.Stopped = false
			cnn.ReceiverConnection.SetOnMessage(cnn.OnMessageReceived)
			cnn.SenderConnection.SetOnMessage(cnn.OnMessageReceived)
			cnn.ReceiverConnection.ListenIncomingPackets()
			go cnn.ProcessACKS()
			break
		}
	}
}

func (cnn *Connection) Stop(){
	cnn.Stopped = true
}
