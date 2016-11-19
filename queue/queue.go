package queue

import . "./message"
import "container/heap"

type PriorityQueue []*Message

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (q PriorityQueue) Push(x interface{}){
	n := len(*q)
	item := x.(*Message)
	item.index = n
	*q = append(*q, item)
}

func (pq *PriorityQueue) Pop() interface{}{
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0: n-1]
	return item
}

func (pq *PriorityQueue) update(msg *Message, msgtext string, priority int){
    msg.msgtext = msgtext
    msg.priority = priority
    heap.Fix(pq, msg.index)
}

func (pq PriorityQueue) Swap(i, j int){
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (q Queue) queueSize() int{
	return len(q)
}