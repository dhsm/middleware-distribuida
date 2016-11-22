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
  //Call connection to send
  
}

func (tpublisher *TopicPublisher) GetTopic() Topic{
  return tpublisher.MyTopic
}
