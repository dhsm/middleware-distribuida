package topic_subscriber

import . "../topic"

type TopicSubscriber struct{
  MyTopic Topic
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic){
  tsubscriber.MyTopic = topic
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}
