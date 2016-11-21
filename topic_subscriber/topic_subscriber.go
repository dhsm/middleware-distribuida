package topic_subscriber

type TopicSubscriber struct{
  MyTopic Topic
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic){
  tsubscriber.MyTopic = topic
}

func (tsubscriber *TopicSubscriber) GetTopic(){
  return tsubscriber.MyTopic
}
