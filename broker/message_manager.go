package broker

import "sync"
import "time"

import . "../packet"

type MessageManager struct{
	Lock *sync.Mutex
	Server Server
	ClientID string
}

func (mm *MessageManager) NewMM(server Server, clientID string, lock *sync.Mutex){
	mm.Lock = lock
	mm.ClientID = clientID
	mm.Server = server
}

func (mm *MessageManager) Execute(){
	for{
		mm.Lock.Lock()
		ref := mm.Server.Receivers[mm.ClientID]
		waitingACK := ref.GetWaitingAck()
		key, msg, found := waitingACK.Pool()
		now := int32(time.Now().Unix())
		if (!found){
			continue
		}else if(now > msg.TimeStamp){
			msg.TimeStamp = now + (5 * 1000)
			msg.Message.Redelivery = true
			temp := mm.Server.Receivers[mm.ClientID]
			channel := temp.ToSend
			pkt_ := Packet{}
			params := []string{mm.ClientID}
			pkt_.CreatePacket(MESSAGE, 0, params, msg.Message)

			channel <- pkt_
		}else{
			amount_sleep := msg.TimeStamp - now
			waitingACK.Add(key, msg)
			time.Sleep(time.Duration(amount_sleep))
		}

		mm.Lock.Unlock()
	}
}