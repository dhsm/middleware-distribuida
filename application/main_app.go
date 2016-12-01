package main

import . "../middleware"

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

  publisher.Send(session.CreateMessage("Pau que nasce torto\nhttp://google.com.br", "perfumaria",1,"m1"))
  publisher.Send(session.CreateMessage("Menina quando...\nhttp://google.com.br", "perfumaria",1,"m2"))
  publisher2.Send(session.CreateMessage("novo macbook\nhttp://google.com.br", "eletronicos",1,"m3"))
  //publisher2.Send(session.CreateMessage("Nunca de endireita","arborismo", 5,"m2"))

  block := make(chan int)

  go notificationListener(&subscriber3)
  go readNotifications(&publisher2, &session, "eletronicos", block)

  <-block

}
