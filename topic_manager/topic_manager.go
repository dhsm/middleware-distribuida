package topic_manager

import . "../message"
import . "../topic"

type TopicManager struct {
  ActiveTopics map[string]*Topic
}

func (tpcManager *TopicManager) addMessageToTopic(topicname string, msg Message){
  if(! topicExists(topicname, tpcManager.ActiveTopics)){
    createTopic(topicname)
  }
  topic := getTopic(topicname, tpcManager.ActiveTopics)
  tpcManager.ActiveTopics[string] = topic
  topic.AddMessage(msg)
}

func (tpcManager *TopicManager) Subscribe(topicname string, clientId string){
  if(! topicExists(topicname, tpcManager.ActiveTopics)){
    createTopic(topicname)
  }
  topic := getTopic(topicname, tpcManager.ActiveTopics)
  topic.AddSubscriber(clientId)
}

func (tpcManager *TopicManager) Unsubscribe(topicname string, clientId string){
  if(topicExists(topicname, tpcManager.ActiveTopics)){
    topic := getTopic(topicname,, tpcManager.ActiveTopics)
    topic.RemoveSubscriber(clientId)
  }
}

func topicExists(topicname string, activeTopics map[string]*Topic) bool{
  val, ok := activeTopics[topicname]
  return ok
}

func createTopic(topicname string){
  topic := Topic{}
  topic.CreateTopic(topicname)
}

func getTopic(topicname string, activeTopics map[string]*Topic){
  val, ok := activeTopics[topicname]
  return val
}
