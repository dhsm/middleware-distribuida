package topic_session

// import . "../topic"
// import . "../message"
// import . "../topic_publisher"
// import . "../topic_subscriber"
// import . "../connection_consumer"

type TopicSession struct {
  // SubscribedList map[string][]MessageListener
  // MyConnectionSendMessage Connection
  // MyMessageListener OutputStream
}

// func (tsession *TopicSession) CreatePublisher(tpc Topic) TopicPublisher {
//   return tsession.createPublisherInternal(tpc).(TopicPublisher)
// }

// func (tsession *TopicSession) createPublisherInternal(tpc Topic) interface{} {
//   tpublisher := TopicPublisher{}
//   tpublisher.CreateTopicPublisher(tpc)

//   return tpublisher
// }

// func (tsession *TopicSession) CreateSubscriber(tpc Topic) TopicSubscriber {
//   return tsession.createSubscriberInternal(tpc).(TopicSubscriber)
// }

// func (tsession *TopicSession) createSubscriberInternal(tpc Topic) interface{} {
//   tsubscriber := TopicSubscriber{}
//   tsubscriber.CreateTopicSubscriber(tpc)

//   tsession.MyConnectionSendMessage.Subscribe(tpc, tsession)
//   //TODO implement HashMap access
//   list_of_subscribers_of_this_topic := tsession.SubscibedList.get(tpc.GetName)
//   list_of_subscribers_of_this_topic.append(tpublisher)
//   tsession.SubscribedList.append(tpc.GetName(),list_of_subscribers_of_this_topic)
  
//   return tsubscriber
// }

// func (tsession *TopicSession) CreateTopic(topicname string) Topic{
//   topic := Topic{}
//   topic.CreateTopic(topicname)
//   return topic
// }

// func (tsession *TopicSession) CreateMessage(msgtext string, priority int) Message{
//   msg := Message{}
//   msg.CreateMessage(msgtext, priority)
//   return msg
// }
