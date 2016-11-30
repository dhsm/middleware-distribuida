package broker

import "net"
import "fmt"
import "log"
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
  server.NextHandlerId = 0
  ln, _ := net.Listen("tcp", port)

  server.Listener = ln
  println("==> Server created!")
}

func (server *Server) Init() {
  for{
    println("Waiting for a new Connection...")
    conn, _ := server.Listener.Accept()
    log.Print("new connetion from ", conn.RemoteAddr())
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
  println("*** Server handle[SENDER]")
  server.Senders[pkt.GetClientID()] = server.Handlers[id]
}

func (server *Server) HandleRegisterReceiver(pkt Packet, id int){
  println("*** Server handle[RECEIVER]")
  println("*** Server receiver added : ",pkt.GetClientID())
  server.Receivers[pkt.GetClientID()] = server.Handlers[id]
}

func (server *Server) HandleSubscribe(pkt Packet){
  println("*** Server handle[SUBSCRIBE]")
  server.MyTopicManager.Subscribe(pkt.GetMessage().Destination, pkt.GetClientID())
}

func (server *Server) HandleUnsubscribe(pkt Packet){
  println("*** Server handle[UNSUBSCRIBE]")
  server.MyTopicManager.Unsubscribe(pkt.GetMessage().Destination, pkt.GetClientID())
}

func (server *Server) HandleCreateTopic(pkt Packet){
  println("*** Server handle[CREATETOPIC]")
  server.MyTopicManager.CreateTopic(pkt.GetMessage().Destination)
}

func (server *Server) HandleMessage(pkt Packet){
  println("*** Server handle[MESSAGE]")
  topic := pkt.GetMessage().Destination
  //server.MyTopicManager.AddMessageToTopic(topic, pkt.GetMessage())
  server.MyTopicManager.AddMessageToTopic(topic, pkt)

  pkt_ := Packet{}
  fmt.Println("*** Server packet ::: ",pkt)
  fmt.Println("*** Server clientID ::: ",pkt.GetClientID())
	params := []string{pkt.GetClientID()}
	pkt_.CreatePacket(ACK.Ordinal(), 0, params, Message{})

  //server.Receivers[pkt.GetClientID()].ToSend <- pkt_
  fmt.Println("*** Server receivers ::: ",server.Receivers)
  // current_receiver := ""
  for key, _ := range server.Receivers {
    fmt.Println("Key:", key)
    // current_receiver = key
  }
  handler := server.Receivers[pkt.GetClientID()]
  handler.ToSend <- pkt_
  // server.Receivers[current_receiver].ToSend <- pkt_
}
