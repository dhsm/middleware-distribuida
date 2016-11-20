package main

import . "./topic"
import . "./message"

func main(){

  msg := Message{}
  msg.CreateMessage("mensagem sobre o assunto",7)

  topic := Topic{}
  topic.CreateTopic("assunto")
  topic.AddMessage(msg)
  println(topic.Messages.Pop())
  topic.AddSubscriber("fulana")
  println(topic.Subscribed[0])
}
