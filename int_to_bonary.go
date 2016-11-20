package main

import "fmt"
import "encoding/binary"
import "bytes"

func main() {
    j := int32(5247)
    fmt.Println(j)
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.BigEndian, j)
    fmt.Println(buf)
    if err != nil {
        fmt.Println(err)
        return
    }
    var k int32
    err = binary.Read(buf, binary.BigEndian, &k)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(k)
}
