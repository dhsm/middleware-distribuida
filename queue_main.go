package main

import . "./queue"

func main(){
    msg := Message{"husadhusaid"}
    pq := make(PriorityQueue, 1)
    pq[0] = msg
    heap.Init(&pq)
}