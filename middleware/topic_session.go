package middleware

import "sync"
//import . "../topic"
import . "../message"
import . "../message_listener"
//import . "../topic_publisher"
//import . "../topic_subscriber"
//import . "../connection"

type TopicSession struct {
  SubscribedList map[string][]TopicSubscriber
  MyConnectionSendMessage Connection
  MyMessageListener MessageListener
}

func (tsession *TopicSession) CreateSession(conn Connection) {
  tsession.SubscribedList = make(map[string][]TopicSubscriber)
  tsession.MyConnectionSendMessage = conn
}

func (tsession *TopicSession) CreateTopicPublisher(tpc Topic) TopicPublisher {
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
  tsubscriber.CreateTopicSubscriber(tpc, tsession)
  //TODO call subscribe
  //tsession.MyConnectionSendMessage.Subscribe(tpc, tsession)

  var mu sync.Mutex

  mu.Lock()
    subscribed_list := tsession.SubscribedList
    list_of_subscribers_of_this_topic := subscribed_list[tpc.GetTopicName()]

    if list_of_subscribers_of_this_topic == nil {
      list := make([]TopicSubscriber, 0)
      list = append(list, tsubscriber)
      list_of_subscribers_of_this_topic = list
    }else{
      list_of_subscribers_of_this_topic = append(list_of_subscribers_of_this_topic, tsubscriber)
    }
    tsession.SubscribedList[tpc.GetTopicName()] = list_of_subscribers_of_this_topic
  mu.Unlock()

  return tsubscriber
}

func (tsession *TopicSession) CreateTopic(topicname string) Topic{
  topic := Topic{}
  topic.CreateTopic(topicname)
  //TODO call createtopic
  //tsession.MyConnectionSendMessage.CreateTopic(topic)
  return topic
}

//TODO check if here we really need to create an empty message
func (tsession *TopicSession) CreateMessage(msgtext string, destination string, priority int, messageid string) Message{
  msg := Message{}
  msg.CreateMessage(msgtext, destination, priority, messageid)
  return msg
}

func (tsession *TopicSession) OnMessageReceived(msg Message) {
  //TODO check if msg really is the object that has SessionAck
  // msg.SetSessionAck(tsession)
  // topic := msg.GetTopic()
  // var mu sync.Mutex
  // mu.Lock()
  //   subscribed_list := tsession.SubscribedList[topic.GetTopicName()]
  //   for _, subscriber := range subscribed_list {
  //       subscriber.onMessage(msg)
  //   }
  // mu.Unlock()
}

func (tsession *TopicSession) Send(msg Message) {
  //TODO call send
  //tsession.MyConnectionSendMessage.Send(msg)
}

func (tsession *TopicSession) CloseSubscriber(publisher TopicSubscriber){
  topic := publisher.GetTopic()
  var mu sync.Mutex
  mu.Lock()
    subscribers := tsession.SubscribedList[topic.GetTopicName()]
    for i, v := range subscribers {
      if(v == publisher ){
        copy(subscribers[i:], subscribers[i+1:])
        subscribers = subscribers[:len(subscribers)-1]
      }
    }

    if len(subscribers) == 0 {
      tsession.Unsubscribe(topic.GetTopicName())
    }
}

func (tsession *TopicSession) Unsubscribe(topic_name string) {
  //TODO call unsubscribe
  //tsession.MyConnectionSendMessage.Unsubscribe(topic_name)
}
