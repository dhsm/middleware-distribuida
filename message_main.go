package main

import . "./message"
import "fmt"
import "encoding/json"

func main(){
  msg := Message{"oi brasil"}

  fmt.Println(msg)
  //If all is well, err will be nil and b will be a []byte containing this JSON data
  msgMarshaled, _ := json.Marshal(msg)
  fmt.Println(msgMarshaled)

  var msgUnmarshaled Message
  _ = json.Unmarshal(msgMarshaled, &msgUnmarshaled)
  fmt.Println(msgUnmarshaled)
}
