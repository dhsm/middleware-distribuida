package message

import "fmt"
import "testing"

func main(){
  msg := Message{"oi brasil"}

  fmt.Println(msg.Msgtext)
}

func TestStruct( t *testing.T){
  var m Message
  m = Message{"1 2 3 testando"}
  if m.Msgtext != "1 2 3 testando"{
    t.Error("Expected string, got ", m)
  }
}
