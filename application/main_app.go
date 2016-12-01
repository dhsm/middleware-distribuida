package main

import . "../middleware"

import "time"

func main() {
  conn := Connection{}
  conn.CreateConnection("127.0.0.1", "8082", "tcp")
  conn.Start()

  session := TopicSession{}
  session.CreateSession(conn)

  topic := session.CreateTopic("perfumaria")
  topic2 := session.CreateTopic("eletronicos")

  publisher := session.CreateTopicPublisher(topic)
  publisher2 := session.CreateTopicPublisher(topic2)
  //publisher2 := session.CreateTopicPublisher(topic)
  subscriber1 := session.CreateTopicSubscriber(topic)
  subscriber2 := session.CreateTopicSubscriber(topic2)
  subscriber3 := session.CreateTopicSubscriber(topic)
  subscriber3 = session.CreateTopicSubscriber(topic2)
  // subscriber2 := session.CreateTopicSubscriber(topic)
  // subscriber3 := session.CreateTopicSubscriber(topic)



  subscriber1.GetTopic()
  subscriber2.GetTopic()
  subscriber3.GetTopic()

  publisher.Send(session.CreateMessage("Pau que nasce torto", "perfumaria",1,"m1"))
  publisher.Send(session.CreateMessage("Menina quando...", "perfumaria",1,"m2"))
  publisher2.Send(session.CreateMessage("novo macbook", "eletronicos",1,"m2"))
  //publisher2.Send(session.CreateMessage("Nunca de endireita","arborismo", 5,"m2"))

  messageFromInput := sendMessage()
  publisher2.Send(session.CreateMessage(messageFromInput, "eletronicos",1,"m2"))

  time.Sleep(time.Second * 3000)

}
