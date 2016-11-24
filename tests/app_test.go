package main

import "testing"
import . "../middleware"

func TestCreateTopic(t *testing.T) {
  var topic Topic
  tsession := TopicSession{}
  topic = tsession.CreateTopic("arvores_topic")
  if topic.GetTopicName() != "arvores_topic" {
    t.Error("Expected arvores_topic, got ",topic.GetTopicName())
  }
}

func TestPublishMessage(t *testing.T) {
  var topic Topic
  tsession := TopicSession{}
  topic = tsession.CreateTopic("arvores_topic")

  tpublisher := tsession.CreateTopicPublisher(topic)

  msg := tsession.CreateMessage("Pau que nasce torto", "arborismo",1,"m1")

  tpublisher.Publish(msg)

  //If message was published, the queue attribute of Connection will store the message
  //TODO check if Connection.queue has the message
}
