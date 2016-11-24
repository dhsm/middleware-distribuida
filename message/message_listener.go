package message

type MessageListener interface {
	OnMessage(msg *Message)
}
