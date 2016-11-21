package main

import "testing"
import . "../topic"
import . "../message"
import . "../topic_session"
//import . "../topic_publisher"

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
  if msgtext != "oi testes"{
    t.Error("Excpected oi testes, got ",msgtext)
  }
}

func TestAddMessageToTopic(t *testing.T) {
  msg := Message{}
  msg.CreateMessage("mensagem sobre o assunto",7)

  topic := Topic{}
  topic.CreateTopic("assunto")
  topic.AddMessage(msg)
  msg_do_topico := topic.Messages.Pop().(*Message)
  msgtext := msg_do_topico.GetText()
  if msgtext != "mensagem sobre o assunto"{
    t.Error("Excpected mensagem sobre o assunto, got ",msgtext)
  }
}

// func TestCreatePublisher(t *testing.T) {
//   var topic Topic
//   tsession := TopicSession{}
//   topic = tsession.CreateTopic("meu_topico_massa")
//   tpublisher := tsession.CreatePublisher(topic)
//
//   msg := Message{}
//   msg.CreateMessage("oi testes", 99)
//
//   tpublisher.Publish(msg)
//
//   topic_with_message = tpublisher.GetTopic()
// }
