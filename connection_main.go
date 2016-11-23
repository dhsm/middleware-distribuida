package main

import "fmt"

import . "./connection"
import . "./packet"

type onPacketReceived func(pkt Packet)

func printReceivedPackets(pkt Packet){
  fmt.Println(pkt)
}

func printReceivedPackets1(pkt Packet){
  fmt.Println(pkt)
}

func printReceivedPackets2(pkt Packet){
  fmt.Println(pkt)
}

func printReceivedPackets3(pkt Packet){
  fmt.Println(pkt)
}

func main() {
	cnn := Connection{}
	cnn.CreateConnection("127.0.0.1", "8081", "tcp")
	fmt.Println(cnn)
	s := make(map[string][]onPacketReceived)
	_, f := s["Soares"]
	if (!f) {
		s["Soares"] = make([]onPacketReceived,0)
	}
	s["Soares"] = append(s["Soares"], printReceivedPackets)
	s["Soares"] = append(s["Soares"], printReceivedPackets1)
	s["Soares"] = append(s["Soares"], printReceivedPackets2)
	s["Soares"] = append(s["Soares"], printReceivedPackets3)

	fmt.Println(s)
}