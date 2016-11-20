package queue

import . "../message"
// import "container/heap"

type PriorityQueue []*Message



func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int){
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq PriorityQueue) Push(x interface{}){
	item := x.(*Message)
	pq = append(pq, item)
}

func (pq PriorityQueue) Pop() interface{}{
	old := pq
	n := len(old)
	item := old[n-1]
	pq = old[0: n-1]
	return item
}