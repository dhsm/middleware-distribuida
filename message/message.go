package message

import "time"

//Good source to learn why the variable name has a capital letter https://www.goinggo.net/2014/03/exportedunexported-identifiers-in-go.html

type Message struct {
  Msgtext string
  Priority int
  Index int //This is necessary because we are using a PriorityQueue
  TimeStamp int32
}

func (msg *Message) CreateMessage(msgtext string, priority int){
  msg.Msgtext = msgtext
  msg.Priority = priority
  msg.Index = -1
  msg.TimeStamp = int32(time.Now().Unix())
}
