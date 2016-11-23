package topic_publisher

import . "../topic"
import . "../message"
import . "../session"

type TopicPublisher struct {
  MyTopic Topic
  SessionSend Session
}

func (tpublisher *TopicPublisher) CreateTopicPublisher(topic Topic, session Session){
  tpublisher.MyTopic = topic
  tpublisher.SessionSend = session
}

func (tpublisher *TopicPublisher) Publish(msg Message){
  tpublisher.SessionSend.Send(msg)
}

func (tpublisher *TopicPublisher) GetTopic() Topic{
  return tpublisher.MyTopic
}
