// We can use channels to synchronize execution
// across goroutines. Here's an example of using a
// blocking receive to wait for a goroutine to finish.

package main

import "fmt"
import "time"

// This is the function we'll run in a goroutine. The
// `done` channel will be used to notify another
// goroutine that this function's work is done.
func worker(done chan bool, t int) {
    fmt.Print("working...")
    time.Sleep(10 * time.Second)
    fmt.Println("done")

    // Send a value to notify that we're done.
    done <- true
}

func main() {

    // Start a worker goroutine, giving it the channel to
    // notify on.
    var t = make(chan int, 10)
    fmt.Println(t)
    done := make(chan bool, 1)
    go worker(done,1)
    go worker(done,2)
    go worker(done,10)

    // Block until we receive a notification from the
    // worker on the channel.
    l:= <-done
    i:= <-done
    j:= <-done
    // k:= <-done
    fmt.Println(l, i, j)
}
