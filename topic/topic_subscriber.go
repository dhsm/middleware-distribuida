package topic

 import . "../message"

type TopicSubscriber struct{
  MyTopic Topic
  SessionReceive interface{}
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic, session TopicSession){
  tsubscriber.MyTopic = topic
  tsubscriber.SessionReceive = session
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}

func (tsubscriber *TopicSubscriber) OnMessage(msg Message) Message {
  //TODO check if this method really exists here
  return msg
}
