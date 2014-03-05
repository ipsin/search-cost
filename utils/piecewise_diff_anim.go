package main

import "fmt"
import "searchcost"
// import "github.com/ipsin/searchcost"

func main() {
  fmt.Printf("Moo\n")
  t := searchcost.NewLinear(1,1)
  fmt.Printf("%s\n", t.String())

  ar := searchcost.NewPiecewise(1,2,1,5,1,5)
  fmt.Printf("%s\n", ar.String())
}
