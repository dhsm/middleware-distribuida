package main

import "testing"
import . "../topic"
import . "../topic_session"

func TestCreateConnection(t *testing.T) {
  //TODO
  // if topicname != "meu_topico_massa" {
  //   t.Error("Expected meu_topico_massa, got ", topicname)
  // }
}

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

  tpublisher := tsession.CreatePublisher(topic)

  msg := tsession.CreateMessage("oi_brasil",99)

  tpublisher.Publish(msg)

  //If message was published, the queue attribute of Connection will store the message
  //TODO check if Connection.queue has the message
}
