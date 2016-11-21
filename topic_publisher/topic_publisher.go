package topic_publisher

import . "../topic"
import . "../message"

type TopicPublisher struct {
  MyTopic Topic
}

func (tpublisher *TopicPublisher) CreateTopicPublisher(topic Topic){
  tpublisher.MyTopic = topic
}

func (tpublisher *TopicPublisher) Publish(msg Message){
  tpublisher.MyTopic.AddMessage(msg)
}

func (tpublisher *TopicPublisher) GetTopic(){
  return tpublisher.MyTopic
}
