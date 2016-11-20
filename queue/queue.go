package queue

import . "../message"
import "container/heap"

type PriorityQueue []*Message

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq *PriorityQueue) Push(x interface{}){
	n := len(*pq)
	item := x.(*Message)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{}{
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0: n-1]
	return item
}

func (pq *PriorityQueue) update(msg *Message, msgtext string, priority int){
    msg.Msgtext = msgtext
    msg.Priority = priority
    heap.Fix(pq, msg.Index)
}

func (pq PriorityQueue) Swap(i, j int){
    pq[i], pq[j] = pq[j], pq[i]
    // pq[i].Index = i
    // pq[j].Index = j
}


