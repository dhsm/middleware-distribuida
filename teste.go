package main

import "time"
import . "./message"
import "fmt"
// import "math/rand"

// A concurrent prime sieve
// from: http://play.golang.org/p/9U22NfrXeq

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- Message) {
    for i := 0; ; i++ {
        msg := Message{}
        msg.CreateMessage("Haha", i)

        select {
            case ch <- msg:
            default:
                println("Channel is full trying again in 0,5 seconds")
        }
        time.Sleep(time.Millisecond * 50)
    }
}

func Consumer(ch <-chan Message, id int){
    // r := rand.New(rand.NewSource(1000))
    for {
        t := <-ch
        fmt.Println(id, " ", t)
        time.Sleep(time.Second * time.Duration(id))
    }
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
    // for {
    //     i := <-in // Receive value from 'in'.
    //     print(i," ")
    //     if i%prime != 0 {

    //         out <- i // Send 'i' to 'out'.
    //     }
    // }
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
    ch := make(chan Message, 100) // Create a new channel.
    go Generate(ch)      // Launch Generate goroutine.
    go Consumer(ch,3)
    go Consumer(ch,4)
    go Consumer(ch,5)

    time.Sleep(100*time.Second)
    // for i := 0; i < 2; i++ {
    //     prime := <-ch
    //     print(prime, "\n")
    //     ch1 := make(chan int)
    //     go Filter(ch, ch1, prime)
    //     ch = ch1
    // }
}


