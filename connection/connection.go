package connection_consumer

import . "./topic"
import . "./queue"

type Connection struct {
  MessageQueue PriorityQueue
  Name string
  Sender interface{}
  Receiver interface{}
}

func (cconsumer *ConnectionConsumer) CreateConnectionConsumer(subsName string, queue PriorityQueue) {
  cconsumer.MessageQueue = queu
  cconsumer.Name = subsName
}

func (cconsumer *ConnectionConsumer) Subscribe(topic Topic, message_listener TopicSession) {

}
