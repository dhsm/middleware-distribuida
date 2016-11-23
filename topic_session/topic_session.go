package topic_session

import . "../topic"
import . "../message"
import . "../topic_publisher"
import . "../topic_subscriber"
import . "../connection"

type TopicSession struct {
  //TODO check if MessageListener type is really OutputStream
  SubscribedList map[string][]OutputStream
  MyConnectionSendMessage Connection
  MyMessageListener OutputStream
}

func (tsession *TopicSession) CreateSession() {
  tsession.SubscribedList = make(map[string][]OutputStream)
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

  tsession.MyConnectionSendMessage.Subscribe(tpc, tsession)

  subscribed_list := tsession.SubscribedList
  list_of_subscribers_of_this_topic := subscribed_list[tpc.GetName()]

  if list_of_subscribers_of_this_topic == nil {
    var list []OutputStream
    list = make([]OutputStream, tsubscriber)
    list_of_subscribers_of_this_topic = list
  }else{
    list_of_subscribers_of_this_topic.append(tsubscriber)
  }
  tsession.SubscribedList[tpc.GetName()] = list_of_subscribers_of_this_topic

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
