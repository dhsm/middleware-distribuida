package middleware

import . "../message"

type TopicPublisher struct {
  MyTopic Topic
  SessionSend *TopicSession
  MessageListener MessageListener
  ClientID string
}

func (tpublisher *TopicPublisher) CreateTopicPublisher(topic Topic, session *TopicSession){
  tpublisher.MyTopic = topic
  tpublisher.SessionSend = session
}

func (tpublisher *TopicPublisher) Publish(msg Message){
  tpublisher.SessionSend.Send(msg)
}

func (tpublisher *TopicPublisher) GetTopic() Topic{
  return tpublisher.MyTopic
}

func (tpublisher *TopicPublisher) Send(msg Message) {
  tpublisher.SessionSend.Send(msg)
}
