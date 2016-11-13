package main

//import . "./message"
import . "./server_request_handler"
//import "encoding/json"
//import "fmt"

func main(){
  srh := ServerRequestHandler{}
  srh.NewSRH("tcp", "127.0.0.1:8081")

  //bytes := srh.Receive()
  sample := []byte{0,0,0,0,1}
  srh.Send(sample)

  //var msgUnmarshaled Message
  //message := json.Unmarshal(bytes, &msgUnmarshaled)
  //fmt.Println(message)
}
