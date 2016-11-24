package main

import . "./broker"

func main(){
  server := Server{}
  server.CreateServer(":8082")
  server.Init()
}
