package middleware

import "sync"
import "math"
import "time"
import "log"
import "errors"
import "reflect"
 // import "fmt"

import "github.com/nu7hatch/gouuid"

import . "../packet"
import . "../message"
import . "../client_request_handler"

type Subscribed struct{
	Map map[string][]MessageListener
}

func (sd *Subscribed) Init() {
	sd.Map = make(map[string][]MessageListener)
}

func (sd *Subscribed) Get(key string) ([]MessageListener, bool){
	l, found := sd.Map[key]
	return l, found
}

func (sd *Subscribed) Set(key string, listeners []MessageListener){
	sd.Map[key] = listeners
}

func (sd *Subscribed) Add(key string, fu MessageListener){
	_, f := sd.Map[key]
	if (!f) {
		sd.Map[key] = make([]MessageListener,0)
	}

	//Checking if this listener is not already in this list
	for _, fun := range sd.Map[key] {
		if(reflect.DeepEqual(fun, fu)){
			return
		}
	}

	sd.Map[key] = append(sd.Map[key], fu)
}

func (sd *Subscribed) Remove(key string, fu MessageListener) bool{
	f, e := sd.Map[key]

	if(!e){
		return e
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

	Subscribed Subscribed
	WaitingACK WaitingACKSafe

	Sessions []TopicSession

	Stopped bool
	Open bool
	Modified bool

	PacketIDGenerator int
}

func (cnn *Connection) CreateConnection(host_ip string, host_port string, host_protocol string){
	cnn.Lock = sync.Mutex{}
	cnn.MessageSent = sync.Cond{L: &sync.Mutex{}}
	cnn.AckReceived = sync.Cond{L: &sync.Mutex{}}

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

func (cnn *Connection) Close(){
	cnn.SetModified()
	cnn.AckReceived.L.Lock()
	cnn.Open = false

	for cnn.WaitingACK.Len() > 0{
		cnn.AckReceived.Wait()
	}

	cnn.ReceiverConnection.Close()
	cnn.SenderConnection.Close()

	cnn.AckReceived.L.Unlock()
}

func (cnn Connection) CreateSession() TopicSession{
	cnn.SetModified()
	tp := TopicSession{}
	tp.CreateSession(cnn)
	cnn.Sessions = append(cnn.Sessions, tp)
	return tp
}

func (cnn *Connection) SendMessage(msg Message) error{
	err := cnn.IsOpen()
	if(err != nil){
		return err
	}

	cnn.WaitingACK.Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix()), msg.MessageID})

	cnn.MessageSent.L.Lock()
	//Broadcasting that there is new messages waiting for an ACK
	cnn.MessageSent.Broadcast()
	cnn.MessageSent.L.Unlock()

	cnn.SetModified()

	pkt := Packet{}
	pkt.CreatePacket(MESSAGE, cnn.PacketIDGenerator, nil, msg)
	cnn.PacketIDGenerator++
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) SubscribeSessionToDestination(topic Topic, fu MessageListener){
	defer cnn.Lock.Unlock()
	cnn.Lock.Lock()
	cnn.Subscribed.Add(topic.GetTopicName(), fu)
}

func (cnn *Connection) UnsubscribeSessionToDestination(topic Topic, fu MessageListener) bool{
	defer cnn.Lock.Unlock()
	cnn.Lock.Lock()
	return cnn.Subscribed.Remove(topic.GetTopicName(), fu)
}

func (cnn *Connection) Subscribe(topic Topic, fu MessageListener) error{
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

func (cnn *Connection) Unsubscribe(topic Topic, fu MessageListener) error{
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
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) CloseSession(ts TopicSession){
	for k, v := range cnn.Subscribed.Map{
		for i, e := range v{
			if(reflect.DeepEqual(e,ts)){
				v[i] = v[len(v)-1]
				v[len(v)-1] = nil
				v = v[:len(v)-1]
				cnn.Subscribed.Set(k, v )
			}
		}
	}
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
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) ProcessACKS(){
	go func () {
		for{
			println("ProcessACKS")
			err := cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}

			if (cnn.WaitingACK.Len() == 0){
				cnn.MessageSent.L.Lock()

				println("ProcessACKS")
				err := cnn.IsOpen()
				if(err != nil){
					log.Print(err)
					break
				}

				cnn.MessageSent.Wait() //Waiting for messages to be sent before stat to process ACKS again
				cnn.MessageSent.L.Unlock()
			}

			println("ProcessACKS")
			err = cnn.IsOpen()
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
					time.Sleep(time.Microsecond * time.Duration(curr))
				}
			}

			println("ProcessACKS")
			err = cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}
		}
	}()
}

func (cnn Connection) OnPacket(pkt Packet){
	if(!cnn.Stopped){
		if(pkt.IsMessage()){
			msg := pkt.Msg
			destination := msg.Destination
			cnn.Lock.Lock()
			sessions, found := cnn.Subscribed.Get(destination)
			if(found){
				for _, session := range sessions{
					session.OnMessage(msg)
				}
			}else{
				println("No sessions")
			}
			cnn.Lock.Unlock()
		}else if(pkt.IsACK()){
			key := pkt.Params[1]
			cnn.Lock.Lock()
			cnn.WaitingACK.Remove(key)
			cnn.AckReceived.Broadcast()
			cnn.Lock.Unlock()
		}
	}
}

func (cnn Connection) Start(){
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
			cnn.ReceiverConnection.SetConnection(cnn)
			cnn.SenderConnection.SetConnection(cnn)
			cnn.ReceiverConnection.ListenIncomingPackets()
			go cnn.ProcessACKS()
			break
		}
	}
}

func (cnn *Connection) Stop(){
	cnn.Stopped = true
}
