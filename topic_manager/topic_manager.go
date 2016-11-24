package topic_manager

import "sync"
import "container/heap"
import . "../message"
//import . "../topic"

type TopicManager struct {
  ActiveTopics []string
  Root Node
}

func (tpcManager *TopicManager) CreateTopicManager(){
  tpcManager.ActiveTopics = make([]string,0)
}

func (tpcManager *TopicManager) Create(topic_name string){
  node := Node{}
  node.CreateNode(topic_name)
}

func (tpcManager *TopicManager) AddMessageToTopic(topicname string, msg Message){
  var addWindowMutex sync.Mutex
  addWindowMutex.Lock()
  tpcManager.ActiveTopics = append(tpcManager.ActiveTopics, topicname)
  msgs := tpcManager.Root.GetMessages()
  heap.Init(&msgs)
  heap.Push(&msgs, &msg)
  tpcManager.Root.SetMessages(msgs)
  tpcManager.Root.SetTotalMessages(tpcManager.Root.GetTotalMessages() + 1)
  addWindowMutex.Unlock()
}

func (tpcManager *TopicManager) PopMessage() Message{
  return tpcManager.popMessage().(Message)
}

func (tpcManager *TopicManager) popMessage() interface{} {
  var addWindowMutex sync.Mutex
  addWindowMutex.Lock()
  messages := tpcManager.Root.GetMessages()
  msgPop := heap.Pop(&messages).(*Message)
  addWindowMutex.Unlock()
  return msgPop
}

func (tpcManager *TopicManager) Subscribe(topic_name string, clientId string){
  subscribers := tpcManager.Root.GetSubscribed()
  subscribers[topic_name] = 1
  tpcManager.Root.SetSubscribed(subscribers)
}

func (tpcManager *TopicManager) Unsubscribe(topicname string, clientId string){
  subscribers := tpcManager.Root.GetSubscribed()
  client := subscribers[clientId]
  if client == 0 {

  }else{
    delete(subscribers, clientId)
  }
  tpcManager.Root.SetSubscribed(subscribers)
}

func (tpcManager *TopicManager) PopActiveTopic() string {
  var addWindowMutex sync.Mutex
  addWindowMutex.Lock()
  defer addWindowMutex.Unlock()
  l := len(tpcManager.ActiveTopics)
  if l == 0 {
    //TODO fu
  }
  res := tpcManager.ActiveTopics[l-1]
  tpcManager.ActiveTopics = tpcManager.ActiveTopics[:l-1]
  return res
}

func (tpcManager *TopicManager) GetNode() Node {
  return tpcManager.getNode().(Node)
}

func (tpcManager *TopicManager) getNode() interface {} {
  return tpcManager.Root
}
