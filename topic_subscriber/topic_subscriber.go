package topic_subscriber

import . "../topic"
import . "../message"

type TopicSubscriber struct{
  MyTopic Topic
  //MessageListener object
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic){
  tsubscriber.MyTopic = topic
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}

func (tsubscriber *TopicSubscriber) OnMessage() Message {
  //TODO
  //Returns the message
  msg := Message{}
  msg.CreateMessage("oi",99)
  return msg
}
