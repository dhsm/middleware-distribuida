package main

import . "./packet"
import . "./message"
import "fmt"

func main() {
	println("Creating message...")
	msg := Message{}
	// msg.CreateMessage("Isso funciona", 10)
	println("Creating packet...")
	pkt := Packet{}
	// params := []string{"leto", "paul", "teg"}
	pkt.CreatePacket(REGISTER_RECEIVER, 0, nil, msg)
	fmt.Println(pkt)
	pkt.CreatePacket(REGISTER_SENDER, 0, nil, msg)
	fmt.Println(pkt)
	print(REGISTER_RECEIVER == REGISTER_SENDER)
}