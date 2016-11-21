package topic_session

import . "../topic"
import . "../topic_publisher"
import . "../topic_subscriber"

type TopicSession struct {


}

func (tsession *TopicSession) CreatePublisherInternal(tpc Topic) interface{} {
  tpublisher := TopicPublisher{}
  tpublisher.CreateTopicPublisher(tpc)
  return tpublisher
}

func (tsession *TopicSession) CreatePublisher(tpc Topic) TopicPublisher {
  return tsession.CreatePublisherInternal(tpc).(TopicPublisher)
}

func (tsession *TopicSession) CreateSubscriberInternal(tpc Topic) interface{} {
  tsubscriber := TopicSubscriber{}
  tsubscriber.CreateTopicSubscriber(tpc)
  return tsubscriber
}

func (tsession *TopicSession) CreateSubscriber(tpc Topic) TopicSubscriber {
  return tsession.CreateSubscriberInternal(tpc).(TopicSubscriber)
}

func (tsession *TopicSession) CreateTopic(topicname string) Topic{
  topic := Topic{}
  topic.CreateTopic(topicname)
  return topic
}
