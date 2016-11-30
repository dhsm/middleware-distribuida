package message

//Good source to learn why the variable name has a capital letter https://www.goinggo.net/2014/03/exportedunexported-identifiers-in-go.html

type Message struct {
  MsgText string
  Priority int
  Destination string
  MessageID string
  Redelivery bool
}

func (msg *Message) CreateMessage(msgtext string, destination string, priority int, messageid string){
  msg.MsgText = msgtext
  msg.Priority = priority
  msg.Destination = destination
  msg.MessageID = messageid
  msg.Redelivery = false
}

func (msg *Message) GetText() string{
  return msg.MsgText
}

func (msg *Message) GetDestination() string{
  return msg.Destination
}
