package broker

import . "../topic_manager"

type Server struct {
	//TODO checar esse socket aqui
	MyServerSocket ServerSocket

	NextHandlerId int
	Handlers map[int]ConnectionHandler
  Senders map[int]ConnectionHandler
  Receivers map[int]ConnectionHandler
  MyTopicManager TopicManager
  MyAdminManager AdminManager
}

func (server *Server) Init(port string) {
  server.Handlers = make(map[int]ConnectionHandler)
  server.Senders = make(map[int]ConnectionHandler)
  server.Receivers = make(map[int]ConnectionHandler)
  tmanager := TopicManager{}
  tmanager.CreateTopicManager()
  server.MyTopicManager = tmanager
  adminmanager := AdminManager{}
  server.MyAdminManager = adminmanager

  for{
    conn, _ := ln.Accept()
    server.MyServerSocket := conn
    connHandler := ConnectionHandler{}
    connHandler.NewCH(server.NextHandlerId,conn)
    server.Handlers[server.NextHandlerId] = connHandler

  }
}
