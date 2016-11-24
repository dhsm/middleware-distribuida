package middleware

import "fmt"
import . "../message"

type TopicSubscriber struct{
  MyTopic Topic
  SessionReceive *TopicSession
  MyMessageListener MessageListener
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic, session *TopicSession){
  tsubscriber.MyTopic = topic
  tsubscriber.SessionReceive = session
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}

func (tsubscriber *TopicSubscriber) OnMessage(msg Message) Message {
  //TODO check if this method really exists here
  fmt.Println("Subscriber falando ",msg)
  return msg
}

func (tsubscriber *TopicSubscriber) SetMessageListener(mlistener MessageListener){
  tsubscriber.MyMessageListener = mlistener
}
