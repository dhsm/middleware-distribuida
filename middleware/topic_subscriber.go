package middleware

import "fmt"
import . "../message"

type TopicSubscriber struct{
  MyTopic Topic
  SessionReceive *TopicSession
  MessageChannel chan Message
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic, session *TopicSession){
  tsubscriber.MyTopic = topic
  tsubscriber.SessionReceive = session
  tsubscriber.MessageChannel = make(chan Message, 10)
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}

func (tsubscriber TopicSubscriber) OnMessage(msg Message){
  println("¥¥¥¥¥¥ TopicSubscriber[ONMESSAGE]")
  tsubscriber.MessageChannel <- msg
  fmt.Println(msg)
}
