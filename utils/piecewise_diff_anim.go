package main

import "fmt"
import "searchcost"
// import "github.com/ipsin/searchcost"

func main() {
  costs := searchcost.CreatePiecewiseSearchCost()

  for t := 1; t < 40; t++ {
    costs.Grow(t+1)
    r1 := costs.Cost(t).OffsetX(1)
    diff := costs.Cost(t + 1).Subtract(&r1)
    if !diff.Equal(&searchcost.ZERO_PIECEWISE) {
      fmt.Printf("** F(x,%d)-F(x+1,%d)=%s\n", t+1, t, diff.String())
    }
  }
}
