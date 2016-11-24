package main

import "fmt"
import "time"

import . "./message"
import . "./packet"

var c = make(chan Packet, 50)

func main() {
  go worker(1)
  pkt := Packet{}
  for i := 0; i < 100; i++ {
  	pkt.CreatePacket(MESSAGE, uint(i), nil, Message{})
    c <- pkt
    fmt.Println("Adding ", pkt)
  }
}

func worker(id int) {
  for {
    _ = <-c
    time.Sleep(time.Second)
  }
}