package queue_manager

import . "../queue"

type QueueManager struct {
  Host string
  Port int
  Queues map[string]Queue
}
