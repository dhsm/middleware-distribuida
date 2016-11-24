package topic_manager

import "sync"
import . "../message"
import . "../topic"

type TopicManager struct {
  ActiveTopics []string
  Root Node
}

func (tpcManager *TopicManager) CreateTopicManager(){
  tpcManager.ActiveTopics = make([]int,0)
}

func (tpcManager *TopicManager) Create(topic_name string){
  node := Node{}
  Node.CreateNode(topic_name)
}

func (tpcManager *TopicManager) AddMessageToTopic(topicname string, msg Message){
  var addWindowMutex sync.Mutex
  addWindowMutex.Lock()
  tpcManager.ActiveTopics = append(tpcManager.ActiveTopics, topicname)
  msgs := tpcManager.Rode.GetMessages()
  heap.Init(&msgs)
  heap.Push(&msgs, &msg)
  tpcManager.Rode.SetMessages(msgs)
  tpcManager.Root.SetTotalMessages(tpcManager.Root.GetTotalMessages() + 1)
  addWindowMutex.Unlock()
}

func (tpcManager *TopicManager) PopMessage() Message{
  return tpcManager.popMessage().(Message)
}

func (tpcManager *TopicManager) popMessage() interface{} {
  var addWindowMutex sync.Mutex
  addWindowMutex.Lock()
  retrun msgPop := heap.Pop(&tpcManager.Root.GetMessages()).(*Message)
  addWindowMutex.Unlock()
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

func (tpcManager *TopicManager) PopActiveTopic() string {
  var addWindowMutex sync.Mutex
  addWindowMutex.Lock()
  defer addWindowMutex.Unlock()
  l := len(tpcManager.ActiveTopics)
  if l == 0 {
    return 0, errors.New("Empty Stack")
  }
  res := tpcManager.ActiveTopics[l-1]
  tpcManager.ActiveTopics = s.s[:l-1]
  return res, nil
}

func (tpcManager *TopicManager) GetNode() Node {
  return tpcManager.getNode().(Node)
}

func (tpcManager *TopicManager) getNode() interface {} {
  return tpcManager.Root
}
