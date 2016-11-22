package main

func main(){
  for{
    tsubscriber :=
    messages_listener := tsubscriber.GetMessagesListener()

    trigger.On("message-arrived", func() {
    // Do Some Task Here.
    fmt.Println("Done")
  })
  }
}
