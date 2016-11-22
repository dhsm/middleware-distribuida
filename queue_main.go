package main

import . "./queue"
import . "./message"
import "container/heap"
import "fmt"

func main(){
    msg := Message{}
    msg.CreateMessage("husadhusaid",7)
    msg2 := Message{}
    msg3 := Message{}
    msg4 := Message{}
    msg2.CreateMessage("yolo", 37000)
    msg3.CreateMessage("olol", 5)
    msg4.CreateMessage("kkkraiolaser", 96)
    pq := make(PriorityQueue, 1)
    pq[0] = &msg
    //pq[1] = &msg2
    //pq[2] = &msg3
    //pq[3] = &msg4
    heap.Init(&pq)
    //heap.Push(&pq, &msg)
    heap.Push(&pq, &msg2)
    heap.Push(&pq, &msg3)
    heap.Push(&pq, &msg4)
    fmt.Println("precisa ser yolo")
    msgPop := heap.Pop(&pq).(*Message)
    fmt.Println(msgPop.Msgtext)
}
