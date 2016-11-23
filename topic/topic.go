package topic

type Topic struct {
  Name string
}

func (tpc *Topic) CreateTopic(name string){
  tpc.Name = name
}

func (tpc *Topic) GetTopicName() string{
  return tpc.Name
}

func (tpc *Topic) RemoveSubscriber(clientId string){
  //TODO implement for real
}
