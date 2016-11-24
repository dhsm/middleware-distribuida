package topic_manager

import . "../message"
import . "../topic"

type TopicManager struct {
  ActiveTopics map[string]int
  Root Node
}

func (tpcManager *TopicManager) CreateTopicManager(){
  tpcManager.ActiveTopics = make(map[string]int)
}

func (tpcManager *TopicManager) Create(topic_name string){
  node := Node{}
  Node.CreateNode(topic_name)
}

func (tpcManager *TopicManager) AddMessageToTopic(topicname string, msg Message){
  tpcManager.ActiveTopics[topic_name] = 1
  msgs := tpcManager.Rode.GetMessages()
  heap.Init(&msgs)
  heap.Push(&msgs, &msg)
  tpcManager.Rode.SetMessages(msgs)
  tpcManager.Root.SetTotalMessages(tpcManager.Root.GetTotalMessages() + 1)
}

func (tpcManager *TopicManager) Subscribe(topic_name string, clientId string){
  subscribers := tpcManager.Root.GetSubscribed()
  subscribers[topic_name] = 1
  tpcManager.Root.SetSubscribed() = subscribers
}

func (tpcManager *TopicManager) Unsubscribe(topicname string, clientId string){
  subscribers := tpcManager.Root.GetSubscribed()
  for i, v := range subscribers {
    if(v == clientId ){
      copy(subscribers[i:], subscribers[i+1:])
      subscribers = subscribers[:len(subscribers)-1]
    }
  }
  tpcManager.Root.SetSubscribed() = subscribers
}
