package topic_manager

import . "../queue"

type Node struct {
  Name string
	Subscribed map[string]int
  //TODO messages needs to be a priorityBLOCKINGqueue
	Messages PriorityQueue
	TotalMessages int
}

func (node *Node) CreateNode(name string){
  node.Name = name
  node.Subscribed = make(map[string]int)
  node.Messages = make(PriorityQueue, 0)
}

func (node *Node) GetSubscribed() map[string]int {
  return node.Subscribed
}

func (node *Node) SetSubscribed(list map[string]int) {
  node.Subscribed = list
}

func (node *Node) GetMessages() PriorityQueue {
  return node.getMessages().(PriorityQueue)
}

func (node *Node) getMessages() interface{} {
  return node.Messages
}

func (node *Node) SetMessages(msgs PriorityQueue) {
  node.Messages = msgs
}

func (node *Node) SetTotalMessages(total int ){
  node.TotalMessages = total
}

func (node *Node) GetTotalMessages() int {
  return node.TotalMessages
}
