package topic_subscriber

import . "../topic"
// import . "../message"
import . "../topic_session"

type TopicSubscriber struct{
  MyTopic Topic
  SessionReceive Session
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic, session Session){
  tsubscriber.MyTopic = topic
  tsubscriber.SessionReceive = session
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}

func (tsubscriber *TopicSubscriber) OnMessage(msg Message) Message {
  //TODO check if this method really exists here
}
