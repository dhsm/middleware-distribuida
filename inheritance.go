package main

import "fmt"

type Shaper interface {
   Area() int
   getstuct() Shaper
}

type Rectangle struct {
   length, width int
}

func (r Rectangle) Area() int {
   return r.length * r.width
}

func (r Rectangle) getstuct() Shaper{
  return r
}

func (r Rectangle) oi(){
  print("Oi Rectangle")
}

type Square struct {
   side int
}

func (sq Square) Area() int {
   return sq.side * sq.side
}

func (sq Square) getstuct() Shaper{
  return sq
}

func (sq Square) oi(){
  print("Oi Square")
}

func main() {
   r := Rectangle{length:5, width:3}
   q := Square{side:5}
   var s Shaper
   s = r
   s.oi()
   shapesArr := [...]Shaper{r, q, s}

   fmt.Println("Looping through shapes for area ...")
   for n, _ := range shapesArr {
       fmt.Println("Shape details: ", shapesArr[n])
       fmt.Println("Area of this shape is: ", shapesArr[n].Area())
   }
}
