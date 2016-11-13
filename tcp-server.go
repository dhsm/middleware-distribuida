package main

import "net"
import "fmt"
import "bufio"
import "strings"

func main(){
	fmt.Println("Lauching server...")

	ln, _ := net.Listen("tcp", "127.0.0.1:8081")

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))
		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))
	}
}
