package message

type MessageWaitingAck struct{
	Message Message
	TimeStamp int32
	MessageID string
}