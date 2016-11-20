package queue_manager_proxy

import . "../client_request_handler"

type QueueManagerProxy struct {
  Queue_name string
}

func (qmp *QueueManagerProxy) NewQMProxy(queue_name string,){
  qmp.Queue_name = queue_name
}

func (qmp *QueueManagerProxy) send(msg Message){
  crh := ClientRequestHandler{}
  crh.NewCRH("tcp", "127.0.0.1:8081")
  packet := Packet{qmp.Queue_name, msg}
  packetMarshaled, _ := json.Marshal(packet)
  crh.send(packetMarshaled)
}
