package packet

import . "../message"

type Packet struct{
  Header string
  Body Message
}
