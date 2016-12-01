package main

import (
    "bufio"
    "fmt"
    "os"
    "github.com/nu7hatch/gouuid"
    . "../middleware"
)

func readNotifications(publisher *TopicPublisher, session *TopicSession,topic string, channel (chan int)){
	i := 0
	for{
		messageFromInput := sendMessage()
		uuid_, _ := uuid.NewV4()
		publisher.Send(session.CreateMessage(messageFromInput, topic,i,uuid_.String()))
		i++
	}
}

func sendMessage() string{
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter the sale description: ")
  title, _ := reader.ReadString('\n')
  fmt.Print("Enter the sale link: ")
  link, _ := reader.ReadString('\n')
  return title + "\n" + link
}
