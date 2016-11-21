package main

import . "./packet"
import . "./message"
import "fmt"

func main() {
	println("Creating message...")
	msg := Message{}
	msg.CreateMessage("Isso funciona", 10)
	println("Creating packet...")
	pkt := Packet{}
	params := []string{"leto", "paul", "teg"}
	pkt.CreatePacket(MESSAGE, 0, params, msg)
	fmt.Println(pkt)
}