package MessageListener

import . "../message"

type MessageListener interface {
	OnMessage(msg *Message)
}
