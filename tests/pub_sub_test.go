package main

import "testing"
import . "../topic"
import . "../message"
import . "../topic_session"

func TestCreateTopic(t *testing.T) {
  var topic Topic
  tsession := TopicSession{}
  topic = tsession.CreateTopic("meu_topico_massa")
  topicname := topic.GetTopicName()
  if topicname != "meu_topico_massa" {
    t.Error("Expected meu_topico_massa, got ", topicname)
  }
}

func TestCreateMessage(t *testing.T) {
  msg := Message{}
  msg.CreateMessage("oi testes", 99)
  msgtext := msg.GetText()
  if msgtext != "oi testes" {
    t.Error("Excpected oi testes, got ",msgtext)
  }
}

func TestAddMessageToTopic(t *testing.T) {
  msg := Message{}
  msg.CreateMessage("mensagem sobre o assunto",7)

  topic := Topic{}
  topic.CreateTopic("assunto")
  topic.AddMessage(msg)
  msg_do_topico := topic.Messages.PopMessage()
  msgtext := msg_do_topico.GetText()
  if msgtext != "mensagem sobre o assunto" {
    t.Error("Expected mensagem sobre o assunto, got ",msgtext)
  }
}

func TestCreatePublisher(t *testing.T) {
  var topic Topic
  tsession := TopicSession{}
  topic = tsession.CreateTopic("meu_topico_massa")
  tpublisher := tsession.CreatePublisher(topic)

  msg := Message{}
  msg.CreateMessage("oi brasil", 99)

  tpublisher.Publish(msg)

  topic_with_message := tpublisher.GetTopic()
  queue := topic_with_message.GetMessages()
  message := queue.PopMessage()
  if message.GetText() != "oi brasil" {
    t.Error("Expected oi brasil, got", message)
  }
}

func TestCreateSubscriber(t *testing.T) {
  var topic Topic
  tsession := TopicSession{}
  topic = tsession.CreateTopic("massa")
  tpublisher := tsession.CreatePublisher(topic)

  msg := Message{}
  msg.CreateMessage("mensagem da hora", 99)

  tpublisher.Publish(msg)

  topic_with_message := tpublisher.GetTopic()

  tsubscriber := tsession.CreateSubscriber(topic_with_message)

  topic_from_subscriber := tsubscriber.GetTopic()
  queue := topic_from_subscriber.GetMessages()
  message := queue.PopMessage()
  if message.GetText() != "mensagem da hora" {
    t.Error("Expected massa, got ",message)
  }
}

func TestMessagePriority(t *testing.T) {
  msg := Message{}
  msg.CreateMessage("mensagem sobre o assunto",6)

  msg2 := Message{}
  msg2.CreateMessage("mensagem 2",7)

  topic := Topic{}
  topic.CreateTopic("assunto")
  topic.AddMessage(msg)
  topic.AddMessage(msg2)
  msg_do_topico := topic.Messages.PopMessage()
  msgtext := msg_do_topico.GetText()
  if msgtext != "mensagem sobre o assunto" {
    t.Error("Expected mensagem sobre o assunto, got ",msgtext)
  }
}
