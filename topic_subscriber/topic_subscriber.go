package topic_subscriber

import . "../topic"
import . "../message"

type TopicSubscriber struct{
  MyTopic Topic
  MessagesListener PriorityQueue
}

func (tsubscriber *TopicSubscriber) CreateTopicSubscriber(topic Topic){
  tsubscriber.MyTopic = topic
  tsubscriber.MessagesListener = make(PriorityQueue, 0)
  heap.Init(&tsubscriber.MessagesListener)
}

func (tsubscriber *TopicSubscriber) GetTopic() Topic{
  return tsubscriber.MyTopic
}

func (tsubscriber *TopicSubscriber) GetMessagesListener() PriorityQueue {
	return tsubscriber.getMessagesListenerInternal().(PriorityQueue)
}

func (tsubscriber *TopicSubscriber) getMessagesListenerInternal() interface{} {
  return tsubscriber.MessagesListener
}

func (tsubscriber *TopicSubscriber) AddMessage(msg Message){
  heap.Push(&tsubscriber.Messages, &msg)
}
