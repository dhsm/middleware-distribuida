package topic

import . "../message"
import . "../queue"
import "container/heap"

type Topic struct {
  Subscribed []string
  Messages PriorityQueue
  Totalmessages int
  Name string
}

func (tpc *Topic) CreateTopic(name string){
  tpc.Name = name
  tpc.Totalmessages = -1
  tpc.Messages = make(PriorityQueue, 0)
  heap.Init(&tpc.Messages)
}

func (tpc *Topic) AddMessage(msg Message){
  heap.Push(&tpc.Messages, &msg)
}

func (tpc *Topic) AddSubscriber(clientId string){
  //TODO implement for real
  tpc.Subscribed = append(tpc.Subscribed,clientId)
}

func (tpc *Topic) RemoveSubscriber(clientId string){
  //TODO implement for real
}