package main

import . "../middleware"

import "time"

func main() {
  conn := Connection{}
  conn.CreateConnection("127.0.0.1", "8082", "tcp")
  conn.Start()

  session := TopicSession{}
  session.CreateSession(conn)

  topic := session.CreateTopic("arborismo")

  publisher := session.CreateTopicPublisher(topic)
  //publisher2 := session.CreateTopicPublisher(topic)
  //subscriber1 := session.CreateTopicSubscriber(topic)
  // subscriber2 := session.CreateTopicSubscriber(topic)
  // subscriber3 := session.CreateTopicSubscriber(topic)



  //subscriber1.SetMessageListener()

  publisher.Send(session.CreateMessage("Pau que nasce torto", "arborismo",1,"m1"))
  //publisher2.Send(session.CreateMessage("Nunca de endireita","arborismo", 5,"m2"))


  time.Sleep(time.Second * 3000)

}
