package main

import . "./message"
import "fmt"

func main(){
  msg := Message{"oi brasil"}

  fmt.Println(msg.Msgtext)
}
