package topic

import . "../message"
import . "../queue"
import "container/heap"

type Topic struct {
  Messages PriorityQueue
  Name string
}

func (tpc *Topic) CreateTopic(name string){
  tpc.Name = name
  tpc.Messages = make(PriorityQueue, 0)
  heap.Init(&tpc.Messages)
}

func (tpc *Topic) AddMessage(msg Message){
  heap.Push(&tpc.Messages, &msg)
}
