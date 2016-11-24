package broker

import . "../packet"
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
  server.MyServerSocket := ServerSocket{}
  server.Handlers = make(map[int]ConnectionHandler)
  server.Senders = make(map[int]ConnectionHandler)
  server.Receivers = make(map[int]ConnectionHandler)
  tmanager := TopicManager{}
  tmanager.CreateTopicManager()
  server.MyTopicManager = tmanager
  adminmanager := AdminManager{}
  server.MyAdminManager = adminmanager

  for{
    //We don't have subtopics, so we don't need to specify the topic
    //topic := server.MyTopicManager.PopActiveTopic();
		messages := server.MyTopicManager.GetNode().GetMessages()
    msgPop := heap.Pop(&messages).(*Message)
    if(msgPop == 0){
      continue
    }
		subscribed = server.MyTopicManager.GetSubscribed()
    for i, v := range subscribed {
      //TODO pass messages to channel
      //server.Receivers.Get
    }
  }
}

type (server *Server) HandleRegisterSender(pkt Packet, id int){

}

type (server *Server) HandleRegisterReceiver(pkt Packet, id int){
  
}

type (server *Server) HandleSubscribe(pkt Packet){
  
}

type (server *Server) HandleUnsubscribe(pkt Packet){
  
}

type (server *Server) HandleCreateTopic(pkt Packet){
  
}

type (server *Server) HandleMessage(pkt Packet){
  
}