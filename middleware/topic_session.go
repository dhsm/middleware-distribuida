package middleware

import "sync"
import "fmt"
import . "../message"

type TopicSession struct {
  SubscribedList map[string][]TopicSubscriber
  MyConnectionSendMessage Connection
  MyMessageListener MessageListener
  mu sync.Mutex
}

func (tsession *TopicSession) CreateSession(conn Connection) {
  println("==> TopicSession created!")
  tsession.SubscribedList = make(map[string][]TopicSubscriber)
  tsession.MyConnectionSendMessage = conn
}

func (tsession *TopicSession) CreateTopicPublisher(tpc Topic) TopicPublisher {
  println("### TopicSession create[PUBLISHER]")
  return tsession.createPublisherInternal(tpc).(TopicPublisher)
}

func (tsession *TopicSession) createPublisherInternal(tpc Topic) interface{} {
  tpublisher := TopicPublisher{}
  tpublisher.CreateTopicPublisher(tpc, tsession)

  return tpublisher
}

func (tsession *TopicSession) CreateTopicSubscriber(tpc Topic) TopicSubscriber {
  println("### TopicSession create[SUBSCRIBER]")
  return tsession.createSubscriberInternal(tpc).(TopicSubscriber)
}

func (tsession *TopicSession) createSubscriberInternal(tpc Topic) interface{} {
  tsubscriber := TopicSubscriber{}
  tsubscriber.CreateTopicSubscriber(tpc, tsession)
  fmt.Println(tsession.MyMessageListener)
  tsession.MyConnectionSendMessage.Subscribe(tpc, tsubscriber)

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
  println("### TopicSession create[TOPIC]")
  topic := Topic{}
  topic.CreateTopic(topicname)
  tsession.MyConnectionSendMessage.CreateTopic(topic)
  return topic
}

//TODO check if here we really need to create an empty message
func (tsession *TopicSession) CreateMessage(msgtext string, destination string, priority int, messageid string) Message{
  println("### TopicSession create[MESSAGE]")
  msg := Message{}
  msg.CreateMessage(msgtext, destination, priority, messageid)
  return msg
}

func (tsession *TopicSession) OnMessageReceived(msg Message) {
  println("### TopicSession [ON_MESSAGE_RECEIVED]")
  tsession.mu.Lock()
  topic := msg.Destination
  list := tsession.SubscribedList[topic]
  for _,subscriber := range list {
    //fmt.Println(subscriber)
    subscriber.OnMessage(msg)
    // index is the index where we are
    // element is the element from someSlice for where we are
    }


  tsession.mu.Unlock()
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
  println("### TopicSession [SEND]")
  //TODO call send
  tsession.MyConnectionSendMessage.SendMessage(msg)
}

func (tsession *TopicSession) CloseSubscriber(publisher TopicSubscriber){
  println("### TopicSession [CLOSE_SUBSCRIBER]")
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
  println("### TopicSession create[UNSUBSCRIBE]")
  //TODO call unsubscribe
  //tsession.MyConnectionSendMessage.Unsubscribe(topic_name)
}
