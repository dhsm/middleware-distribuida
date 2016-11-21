package topic_session

import . "../topic"
import . "../topic_publisher"

type TopicSession struct {


}

func (tsession *TopicSession) CreatePublisher(tpc Topic) TopicPublisher{
  tpublisher := TopicPublisher{}
  tpublisher.CreateTopicPublisher(tpc)
  return tpublisher
}

func (tsession *TopicSession) CreateSubscriber(tpc Topic){

}

func (tsession *TopicSession) CreateTopic(topicname string) Topic{
  topic := Topic{}
  topic.CreateTopic(topicname)
  return topic
}
