package topic_session

import . "../topic"
import . "../topic_publisher"
import . "../topic_subscriber"

type TopicSession struct {


}

func (tsession *TopicSession) CreatePublisher(tpc Topic) TopicPublisher{
  tpublisher := TopicPublisher{}
  tpublisher.CreateTopicPublisher(tpc)
  return tpublisher
}

func (tsession *TopicSession) CreateSubscriber(tpc Topic) TopicSubscriber{
  tsubscriber := TopicSubscriber{}
  tsubscriber.CreateTopicSubscriber(tpc)
  return tsubscriber
}

func (tsession *TopicSession) CreateTopic(topicname string) Topic{
  topic := Topic{}
  topic.CreateTopic(topicname)
  return topic
}
