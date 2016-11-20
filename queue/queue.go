package queue

import . "../message"
import "container/heap"
import "fmt"

type PriorityQueue []*Message

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
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

func main(){
	// Some items and their priorities.
	items := map[string]int{
		"Banana is yellow": -1, "Apple is rich": 1, "Pear is for tea": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
