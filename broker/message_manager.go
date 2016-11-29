package broker

import "sync"
import "time"

type MessageManager struct{
	Lock sync.Mutex
	Server Server
	ClientID string
}

type (mm *MessageManager) NewMM(server Server, clientID string, lock sync.Mutex){
	mm.Lock = lock
	mm.ClientID = clientID
	mm.Server = server
}

type (mm *MessageManager) Execute(){
	for{
		mm.Lock.Lock()
		ref = server.Receivers[mm.ClientID]
		waitingACK = ref.getWaitingAck()
		msg, found = waitingACK.Peek()
		curr := int32(time.Now().Unix())
		if (!found){
			continue
		}else if(curr > msg.TimeStamp){

		}

		mm.Lock.Unlock()
	}
}
