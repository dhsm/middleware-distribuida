package topic_session

import . "../topic"

type TopicSession struct {


}

func (tsession *TopicSession) createPublisher(tpc Topic){

}

func (tsession *TopicSession) createSubscriber(tpc Topic){

}

func (tsession *TopicSession) createTopic(topicname string) Topic{
  topic := Topic{}
  topic.CreateTopic(topicname)
  return topic
}
