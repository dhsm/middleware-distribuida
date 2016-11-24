package broker

import "net"
import "fmt"
import . "../packet"
import . "../message"
import . "../topic_manager"

type Server struct {
	//TODO checar esse socket aqui
	//MyServerSocket net.Conn
  Listener net.Listener
	NextHandlerId int
	Handlers map[int]*ConnectionHandler
  Senders map[string]*ConnectionHandler
  Receivers map[string]*ConnectionHandler
  MyTopicManager TopicManager
  //MyAdminManager AdminManager
}

func (server *Server) CreateServer(port string) {
  server.Handlers = make(map[int]*ConnectionHandler)
  server.Senders = make(map[string]*ConnectionHandler)
  server.Receivers = make(map[string]*ConnectionHandler)
  tmanager := TopicManager{}
  tmanager.CreateTopicManager()
  server.MyTopicManager = tmanager
  //adminmanager := AdminManager{}
  //server.MyAdminManager = adminmanager
  server.NextHandlerId = 0
  ln, _ := net.Listen("tcp", port)

  server.Listener = ln
}

func (server *Server) Init() {
  for{
    fmt.Println(server.Senders)
    println("waiting connection")
    conn, _ := server.Listener.Accept()
    //server.MyServerSocket = conn
    id := server.getNextInt()
    connHandler := ConnectionHandler{}
    connHandler.NewCH(id,conn,*server)
    server.Handlers[id] = &connHandler
    go connHandler.Execute()
  }
}

func (server *Server) getNextInt() int{
  ret := server.NextHandlerId
  server.NextHandlerId++
  return ret
}

func (server *Server) HandleRegisterSender(pkt Packet, id int){
  server.Senders[pkt.GetClientID()] = server.Handlers[id]
}

func (server *Server) HandleRegisterReceiver(pkt Packet, id int){
  server.Receivers[pkt.GetClientID()] = server.Handlers[id]
}

func (server *Server) HandleSubscribe(pkt Packet){
  server.MyTopicManager.Subscribe(pkt.GetMessage().Destination, pkt.GetClientID())
}

func (server *Server) HandleUnsubscribe(pkt Packet){
  server.MyTopicManager.Unsubscribe(pkt.GetMessage().Destination, pkt.GetClientID())
}

func (server *Server) HandleCreateTopic(pkt Packet){
  server.MyTopicManager.CreateTopic(pkt.GetMessage().Destination)
}

func (server *Server) HandleMessage(pkt Packet){
  topic := pkt.GetMessage().Destination
  server.MyTopicManager.AddMessageToTopic(topic, pkt.GetMessage())

  pkt_ := Packet{}
	params := []string{pkt.GetClientID()}
	pkt_.CreatePacket(ACK, 0, params, Message{})

  server.Receivers[pkt.GetClientID()].ToSend <- pkt_

}
