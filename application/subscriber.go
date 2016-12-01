package main

import (
    "github.com/deckarep/gosx-notifier"
    "strings"
    "log"
    . "../middleware"
)

func notificationListener(subscriber *TopicSubscriber){
  for{
    msg := <-subscriber.MessageChannel
    sendNotification(msg.MsgText)
  }
}

func sendNotification(str string){
	str_spt := strings.Split(str, "\n")

	 //At a minimum specifiy a message to display to end-user.
    note := gosxnotifier.NewNotification("Click here to save!")

    //Optionally, set a title
    note.Title = "Sale Alert!"

    //Optionally, set a subtitle
    note.Subtitle = str_spt[0]

    //Optionally, set a sound from a predefined set.
    note.Sound = gosxnotifier.Glass

    //Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
    note.Group = "bss3shsm.cin.ufpe.br"

    //Optionally, set a sender (Notification will now use the Safari icon)
    note.Sender = "com.apple.Safari"

    //Optionally, specifiy a url or bundleid to open should the notification be
    //clicked.

    println(str_spt[1])
    note.Link = str_spt[1]//or BundleID like: com.apple.Terminal

    //Then, push the notification
    err := note.Push()

    //If necessary, check error
    if err != nil {
        log.Println("Uh oh!")
    }
}
