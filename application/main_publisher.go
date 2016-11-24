package main

import . "../middleware"

func main() {
  conn := Connection{}
  conn.CreateConnection("127.0.0.1", ":8081", "tcp")

  session := TopicSession{}
  session.CreateSession(conn)

  topic := session.CreateTopic("arborismo")

  publisher := session.CreateTopicPublisher(topic)

  conn.Start()
  
  publisher.Send(session.CreateMessage("Pau que nasce torto", "arborismo",1,"m1"))
  publisher.Send(session.CreateMessage("Nunca de endireita","arborismo", 5,"m2"))
}
