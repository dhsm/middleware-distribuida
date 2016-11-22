package topic_session

import . "../topic"
import . "../message"
import . "../topic_publisher"
import . "../topic_subscriber"

type TopicSession struct {


}

func (tsession *TopicSession) CreatePublisher(tpc Topic) TopicPublisher {
  return tsession.createPublisherInternal(tpc).(TopicPublisher)
}

func (tsession *TopicSession) createPublisherInternal(tpc Topic) interface{} {
  tpublisher := TopicPublisher{}
  tpublisher.CreateTopicPublisher(tpc)
  return tpublisher
}

func (tsession *TopicSession) CreateSubscriber(tpc Topic) TopicSubscriber {
  return tsession.createSubscriberInternal(tpc).(TopicSubscriber)
}

func (tsession *TopicSession) createSubscriberInternal(tpc Topic) interface{} {
  tsubscriber := TopicSubscriber{}
  tsubscriber.CreateTopicSubscriber(tpc)
  return tsubscriber
}

func (tsession *TopicSession) CreateTopic(topicname string) Topic{
  topic := Topic{}
  topic.CreateTopic(topicname)
  return topic
}

func (tsession *TopicSession) CreateMessage(msgtext string, priority int) Message{
  msg := Message{}
  msg.CreateMessage(msgtext, priority)
  return msg
}
