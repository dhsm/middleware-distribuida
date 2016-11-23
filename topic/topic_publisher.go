package topic

import . "../message"
// import . "../topic_session"

type TopicPublisher struct {
  MyTopic Topic
  SessionSend TopicSession
}

func (tpublisher *TopicPublisher) CreateTopicPublisher(topic Topic, session TopicSession){
  tpublisher.MyTopic = topic
  tpublisher.SessionSend = session
}

func (tpublisher *TopicPublisher) Publish(msg Message){
  tpublisher.SessionSend.Send(msg)
}

func (tpublisher *TopicPublisher) GetTopic() Topic{
  return tpublisher.MyTopic
}
