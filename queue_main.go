package main

import . "./queue"
import . "./message"
import "container/heap"
import "fmt"

func main(){
    msg := Message{"husadhusaid", 97}
    msg2 := Message{"yolo", 2}
    msg3 := Message{"olol", 55}
    msg4 := Message{"kkkraiolaser", 97}
    pq := make(PriorityQueue, 1)
    pq[0] = &msg
    //pq[1] = &msg2
    //pq[2] = &msg3
    //pq[3] = &msg4
    heap.Init(&pq)
    heap.Push(&pq, &msg)
    heap.Push(&pq, &msg2)
    heap.Push(&pq, &msg3)
    heap.Push(&pq, &msg4)
    fmt.Println("precisa ser husadhusaid")
    msgPop := heap.Pop(&pq).(*Message)
    fmt.Print(*msgPop.Msgtext)
}
