package broker

import "sync"
//import "fmt"
import . "../packet"
import  "errors"
//import . "../broker"
//import . "../topic"

type TopicManager struct {
  TopicMap map[string][]string
  Server *Server
  mu sync.Mutex
}

func (tpcManager *TopicManager) CreateTopicManager(server *Server){
  println("==> TopicManager created!")
  tpcManager.TopicMap = make(map[string][]string)
  tpcManager.Server = server
}

func (tpcManager *TopicManager) CreateTopic(topic_name string){
  defer tpcManager.mu.Unlock()
  tpcManager.mu.Lock()
  println("!!! TopicManager create[TOPIC]")
  tpcManager.TopicMap[topic_name] = []string{}
  //fmt.Println(tpcManager.TopicMap)
}

func (tpcManager *TopicManager) DeleteTopic(topic_name string){
  defer tpcManager.mu.Unlock()
  tpcManager.mu.Lock()
  println("!!! TopicManager delete[TOPIC]")
  delete(tpcManager.TopicMap,topic_name)
}

func (tpcManager *TopicManager) AddMessageToTopic(topic_name string, msg Packet) error{
  println("!!! TopicManager add[MESSAGE_TO_TOPIC]")
  subscribers, found := tpcManager.TopicMap[topic_name]
  if(!found){
    return errors.New("Tried to send a message to a non existing topic")
  }else{
    for _,val := range subscribers {
      handler := tpcManager.Server.Receivers[val]
      handler.ToSend <- msg
    }
  }
  return nil
}

func (tpcManager *TopicManager) Subscribe(topic_name string, clientID string) error{
  defer tpcManager.mu.Unlock()
  tpcManager.mu.Lock()
  println("!!! TopicManager [SUBSCRIBE]")
  subscribers, found := tpcManager.TopicMap[topic_name]
  if(!found){
    //fmt.Println(len(topic_name), "#", clientID)
    //fmt.Println(tpcManager.TopicMap)
    return errors.New("Tried to subscribe to a non existing topic")
  }else{
    subscribers = append(subscribers, clientID)
    tpcManager.TopicMap[topic_name] = subscribers
    //fmt.Println(tpcManager.TopicMap[topic_name])
  }
  return nil
}

func (tpcManager *TopicManager) Unsubscribe(topic_name string, clientID string) error{
  defer tpcManager.mu.Unlock()
  tpcManager.mu.Lock()
  println("!!! TopicManager [UNSUBSCRIBE]")
  subscribers, found := tpcManager.TopicMap[topic_name]
  if(!found){
    return errors.New("Tried to unsubscribe to a non existing topic")
  }else{
    for i,val := range subscribers {
      if(val == clientID){
        subscribers[i] = subscribers[len(subscribers)-1]
        subscribers = subscribers[:len(subscribers)-1]
      }
    }
    tpcManager.TopicMap[topic_name] = subscribers
  }
  return nil
}
