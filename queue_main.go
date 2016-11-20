package main

import . "./queue"
import . "./message"
import "container/heap"

func main(){
    msg := Message{"husadhusaid", 97, 1}
    pq := make(PriorityQueue, 1)
    pq[0] = &msg
    heap.Init(&pq)
}
