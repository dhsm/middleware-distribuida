package queue

import . "../packet"
import "container/heap"

type PriorityQueue []*Packet

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
	item := x.(*Packet)
	item.Index = n
	*pq = append(*pq, item)
	//values, err := trigger.Fire("message-arrived", pq.Pop()) 
}

func (pq *PriorityQueue) Pop() interface{}{
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0: n-1]
	return item
}

func (pq *PriorityQueue) update(pkt *Packet, msgtext string, priority int){
	pkt.MsgText = msgtext
	pkt.Priority = priority
	heap.Fix(pq, pkt.Index)
}
