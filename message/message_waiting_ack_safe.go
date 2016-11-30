package message

import "sync"
// import "fmt"

type WaitingACKSafe struct{
	sync.Mutex
	Map map[string]MessageWaitingAck
}

func (was *WaitingACKSafe) Init() {
	was.Map = make(map[string]MessageWaitingAck)
}

func (was *WaitingACKSafe) Get(key string) (MessageWaitingAck, bool){
	defer was.Unlock()
	was.Lock()
	l, found := was.Map[key]
	return l, found
}

func (was *WaitingACKSafe) Add(key string, msg MessageWaitingAck){
	defer was.Unlock()
	was.Lock()
	was.Map[key] = msg
}

func (was *WaitingACKSafe) Remove(key string){
	defer was.Unlock()
	was.Lock()
	delete(was.Map,key)
	// fmt.Println(was.Map)
}

func (was *WaitingACKSafe) Len() int{
	defer was.Unlock()
	was.Lock()
	return len(was.Map)
}

func (was *WaitingACKSafe) Peek() (string, MessageWaitingAck, bool){
	defer was.Unlock()
	was.Lock()
	for k, e := range was.Map {
		return k, e, true
	}

	return "", MessageWaitingAck{}, false
}

func (was *WaitingACKSafe) Pool() (string, MessageWaitingAck, bool){
	defer was.Unlock()
	was.Lock()

	for k, e := range was.Map {
		delete(was.Map,k)
		return k, e, true
	}

	return "", MessageWaitingAck{}, false
}
