package topic_session

import "sync"
import . "../topic"
import . "../message"
import . "../message_listener"
import . "../topic_publisher"
import . "../topic_subscriber"
import . "../connection"

type TopicSession struct {
  SubscribedList SubscribedSafe
  MyConnectionSendMessage Connection
  MyMessageListener MessageLister
}

func (tsession *TopicSession) CreateSession(conn Connection) {
  tsession.SubscribedList = make(map[string][]OutputStream)
  tsession.MyConnectionSendMessage = conn
}

func (tsession *TopicSession) CreatePublisher(tpc Topic) TopicPublisher {
  return tsession.createPublisherInternal(tpc).(TopicPublisher)
}

func (tsession *TopicSession) createPublisherInternal(tpc Topic) interface{} {
  tpublisher := TopicPublisher{}
  tpublisher.CreateTopicPublisher(tpc, tsession)

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
  tsession.MyConnectionSendMessage.CreateTopic(topic)
  return topic
}

func (tsession *TopicSession) CreateMessage(msgtext string, priority int) Message{
  msg := Message{}
  msg.CreateMessage(msgtext, priority)
  return msg
}

func (tsession *TopicSession) onMessageReceived(msg Message) {
  //TODO check if msg really is the object that has SessionAck
  msg.SetSessionAck(tsession)
  topic := msg.GetTopic()
  var mu sync.Mutex
  mu.Lock()
    subscribed_list := tsession.SubscribedList[topic.GetTopicName()]
    for _, subscriber := range subscribed_list {
        subscriber.onMessage(msg)
    }
  mu.Unlock()
}

func (tsession *TopicSession) Send(msg Message) {
  tsession.MyConnectionSendMessage.send(msg)
}
