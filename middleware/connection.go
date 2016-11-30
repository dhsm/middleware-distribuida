package middleware

import "sync"
import "math"
import "time"
import "log"
import "errors"
import "reflect"
import "fmt"

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
	println("==> Conection created!")
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

func (cnn *Connection) Close(){
	println("+++ Conection [CLOSE]")
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
	println("+++ Conection create[SESSION]")
	cnn.SetModified()
	tp := TopicSession{}
	tp.CreateSession(cnn)
	cnn.Sessions = append(cnn.Sessions, tp)
	return tp
}

func (cnn *Connection) SendMessage(msg Message) error{
	//println("+++ Conection send[MESSAGE]")
	err := cnn.IsOpen()
	//println("+++ Conection WTF")

	if(err != nil){
		fmt.Println(err)

		return err
	}
	//println("+++ Conection add WaitingACK")
	cnn.WaitingACK.Add(msg.MessageID, MessageWaitingAck{msg, int32(time.Now().Unix() + (5 * 1000)), msg.MessageID})

	cnn.MessageSent.L.Lock()
	//Broadcasting that there is new messages waiting for an ACK
	cnn.MessageSent.Broadcast()
	cnn.MessageSent.L.Unlock()

	cnn.SetModified()

	pkt := Packet{}
	params := []string{cnn.ClientID,msg.GetDestination()}
	pkt.CreatePacket(MESSAGE.Ordinal(), cnn.PacketIDGenerator, params, msg)
	//pkt.CreatePacket(MESSAGE, cnn.PacketIDGenerator, nil, msg)
	cnn.PacketIDGenerator++
	//println("+++ Conection SendAsync packet")
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) SubscribeSessionToDestination(topic Topic, fu MessageListener){
	println("+++ Conection subscribe[SESSION_TO_DESTINATION]")
	defer cnn.Lock.Unlock()
	cnn.Lock.Lock()
	cnn.Subscribed.Add(topic.GetTopicName(), fu)
}

func (cnn *Connection) UnsubscribeSessionToDestination(topic Topic, fu MessageListener) bool{
	println("+++ Conection UNsubscribe[SESSION_TO_DESTINATION]")
	defer cnn.Lock.Unlock()
	cnn.Lock.Lock()
	return cnn.Subscribed.Remove(topic.GetTopicName(), fu)
}

func (cnn *Connection) Subscribe(topic Topic, fu MessageListener) error{
	println("+++ Conection [SUBSCRIBE]")
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	cnn.SetModified()
	cnn.SubscribeSessionToDestination(topic, fu)
	pkt := Packet{}
	params := []string{cnn.ClientID, topic.GetTopicName()}
	pkt.CreatePacket(SUBSCRIBE.Ordinal(), cnn.PacketIDGenerator, params, Message{})
	cnn.PacketIDGenerator++
	return cnn.SenderConnection.Send(pkt)
}

func (cnn *Connection) Unsubscribe(topic Topic, fu MessageListener) error{
	println("+++ Conection [UNSUBSCRIBE]")
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
		pkt.CreatePacket(UNSUBSCRIBE.Ordinal(), cnn.PacketIDGenerator, params, Message{})
		cnn.PacketIDGenerator++
		return cnn.SenderConnection.Send(pkt)
	}

	return nil
}

func (cnn *Connection) AcknowledgeMessage(msg Message, ts TopicSession) error{
	println("+++ Conection [ACK_MESSAGE]")
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	cnn.SetModified()
	pkt := Packet{}
	params := []string{cnn.ClientID, msg.MessageID}
	pkt.CreatePacket(ACK.Ordinal(), cnn.PacketIDGenerator, params, Message{})
	cnn.PacketIDGenerator++
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) CloseSession(ts TopicSession){
	println("+++ Conection [CLOSE_SESSION]")
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
	println("+++ Conection create[TOPIC]")
	err := cnn.IsOpen()
	if(err != nil){
		log.Print(err)
		return err
	}
	pkt := Packet{}
	params := []string{cnn.ClientID, tp.GetTopicName()}
	pkt.CreatePacket(CREATE_TOPIC.Ordinal(), cnn.PacketIDGenerator, params, Message{})
	cnn.PacketIDGenerator++
	cnn.SenderConnection.SendAsync(pkt)
	return nil
}

func (cnn *Connection) ProcessACKS(){
	println("+++ Conection process[ACKS]")
	go func() {
		for{
			//println("ProcessACKS")
			err := cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}

			if (cnn.WaitingACK.Len() == 0){
				println("Sem ACKS")
				cnn.MessageSent.L.Lock()

				//println("ProcessACKS")
				err := cnn.IsOpen()
				if(err != nil){
					log.Print(err)
					break
				}

				cnn.MessageSent.Wait() //Waiting for messages to be sent before stat to process ACKS again
				cnn.MessageSent.L.Unlock()
			}else{
				println(cnn.WaitingACK.Len())
			}

			//println("ProcessACKS")
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
					println("RESENDING MESSAGE TIMEDOUT!")
					cnn.SendMessage(maa.Message)
				}else{
					time.Sleep(time.Microsecond * time.Duration(curr))
				}
			}

			
			err = cnn.IsOpen()
			if(err != nil){
				log.Print(err)
				break
			}
		}
	}()
}

func (cnn *Connection) OnPacket(pkt Packet){
	println("+++ Conection [ON_PACKET]")
	if(!cnn.Stopped){
		if(pkt.IsMessage()){
			msg := pkt.GetMessage()
			destination := msg.Destination
			cnn.Lock.Lock()
			sessions, found := cnn.Subscribed.Get(destination)
			fmt.Println(sessions)
			if(found){
				for _, session := range sessions{
					session.OnMessage(msg)
					println("chamou on message de", destination)
				}
			}else{
				println("No sessions")
			}
			cnn.Lock.Unlock()
		}else if(pkt.IsACK()){
			fmt.Println("Precessing packet [OnPacket] ",pkt)
			fmt.Println("Length of params array on OnPacket ",len(pkt.Params))
			if(len(pkt.Params) < 1){
				fmt.Println(errors.New("Params (an slice) of packet has no ACK index"))
				//panic(fmt.Sprintf("halp"))
			}else{
				key := pkt.Params[1]
				// println("A porra da KEY tinha de ser m1 ", key)
				// time.Sleep(time.Second * 10)
				cnn.Lock.Lock()
				cnn.WaitingACK.Remove(key)
				cnn.AckReceived.Broadcast()
				cnn.Lock.Unlock()
			}
		}
	}else{
		println("Connection stopped!")
	}
}

func (cnn *Connection) Start(){
	println("+++ Conection [START]")
	if(!cnn.Open){
		tries := 0
		for tries <= 5 {
			tries++
			errr := cnn.ReceiverConnection.NewCRH(cnn.HostProtocol, cnn.HostIp, cnn.HostPort, true, cnn.GetClientID())
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
			println("+++ Conection [ReceiverConnection]SetConnection")
			cnn.ReceiverConnection.SetConnection(cnn)
			println("+++ Conection [SenderConnection]SetConnection")
			cnn.SenderConnection.SetConnection(cnn)
			println("+++ Conection [ReceiverConnection]ListenIncomingPackets")
			cnn.ReceiverConnection.ListenIncomingPackets()
			go cnn.ProcessACKS()
			break
		}
	}
}

func (cnn *Connection) Stop(){
	println("+++ Conection [STOP]")
	cnn.Stopped = true
}
