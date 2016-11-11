package message

import "fmt"
import "testing"

func TestStruct( t *testing.T){
  var m Message
  m = Message{"1 2 3 testando"}
  if m.Msgtext != "1 2 3 testando"{
    t.Error("Expected string, got ", m)
  }
}
