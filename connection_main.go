package main

import "fmt"
// import "reflect"

import . "./middleware"
import . "./packet"

type onPacketReceived func(pkt Packet)

type Tst struct{
	k int
}

func (t *Tst) printReceivedPackets(pkt Packet){
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

	// t1 := Tst{1}
	// t2 := Tst{2}

	// println(reflect.DeepEqual(t1,t2))

	// s := make(map[string][]onPacketReceived)
	// _, f := s["Soares"]
	// if (!f) {
	// 	s["Soares"] = make([]onPacketReceived,0)
	// }

	// for _, fun := range s["Soares"] {
	// 	if(reflect.DeepEqual(fun, t1.printReceivedPackets)){
	// 		println("1")
	// 		return
	// 	}
	// }

	// s["Soares"] = append(s["Soares"], t2.printReceivedPackets)

	// for _, fun := range s["Soares"] {
	// 	if(reflect.DeepEqual(fun, t2.printReceivedPackets)){
	// 		println("Fudeu!	")
	// 		return
	// 	}
	// }

	// s["Soares"] = append(s["Soares"], t2.printReceivedPackets)
	// // s["Soares"] = append(s["Soares"], printReceivedPackets2)

	// for _, fun := range s["Soares"] {
	// 	if(reflect.ValueOf(fun).Pointer() == reflect.ValueOf(printReceivedPackets2).Pointer()){
	// 		return
	// 	}
	// }

	// s["Soares"] = append(s["Soares"], printReceivedPackets2)

	// for _, fun := range s["Soares"] {
	// 	if(reflect.ValueOf(fun).Pointer() == reflect.ValueOf(printReceivedPackets3).Pointer()){
	// 		return
	// 	}
	// }

	// s["Soares"] = append(s["Soares"], printReceivedPackets3)

	// for i, elem := range s["Soares"] {
	// 	if(reflect.ValueOf(elem).Pointer() == reflect.ValueOf(printReceivedPackets2).Pointer()){
	// 		println("Found")
	// 		s["Soares"][i] = s["Soares"][len(s["Soares"])-1]
	// 		s["Soares"][len(s["Soares"])-1] = nil
	// 		s["Soares"] = s["Soares"][:len(s["Soares"])-1]
	// 		break
	// 	}
	// }

	// fmt.Println(s)
}