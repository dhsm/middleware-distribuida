package main

import "net"
import "fmt"
import "bufio"
import "os"

func main(){
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	//fmt.Print(err)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Test to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprint(conn, text + "\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
}
