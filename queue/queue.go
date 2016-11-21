package queue

import . "../message"
import "container/heap"

type PriorityQueue []*Message

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	x,y := int32(pq[i].Priority),int32(pq[j].Priority)
	x_t, y_t := pq[i].TimeStamp,pq[j].TimeStamp

	if(x -1 == 0){
		x = -1
	}
	if(y -1 == 0){
		y = -1
	}
	return x_t/x < y_t/y
}

func (pq PriorityQueue) Swap(i, j int){
    pq[i], pq[j] = pq[j], pq[i]
		pq[i].Index = i
		pq[j].Index = j
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
	item.Index = -1
	*pq = old[0: n-1]
	return item
}

func (pq *PriorityQueue) PopMessage() *Message {
	return pq.Pop().(*Message)
}

func (pq *PriorityQueue) update(msg *Message, msgtext string, priority int){
	msg.Msgtext = msgtext
	msg.Priority = priority
	heap.Fix(pq, msg.Index)
}
