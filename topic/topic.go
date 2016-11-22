package topic

import . "../message"
import . "../queue"
import "container/heap"

type Topic struct {
  Name string
}

func (tpc *Topic) CreateTopic(name string){
  tpc.Name = name
}

func (tpc *Topic) GetTopicName() string{
  return tpc.Name
}
