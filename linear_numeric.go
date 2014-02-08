package searchcost

import "math"
import "sync"

// Key, this is the cost of searching the range x,x+1,...,x+n
type LinearSearchRange struct {
  x, n int
}

type LinearSearchResult struct {
  cost uint64
  minSplitPoints []int
}

var zeroCost = LinearSearchResult{0, []int{}}

func CalculateNumericRange(r LinearSearchRange,
  results *map[LinearSearchRange]LinearSearchResult,
  mutex *sync.Mutex) LinearSearchResult {

  switch {
  case r.n == 0:
    return zeroCost
  case r.n == 1:
    return LinearSearchResult{uint64(r.x), []int{0}}
  case r.n == 2:
    return LinearSearchResult{uint64(r.x + 1), []int{1}}
  } 

  var minCost uint64 = math.MaxUint64
  var cost uint64
  var minSplitPoints []int = []int{}

  for k := 1; k < r.n; k++ {
    mutex.Lock()
    v, has := (*results)[r]
    mutex.Unlock()
    if has {
      return v
    }
   
    leftCost := CalculateNumericRange(LinearSearchRange{r.x, k-1}, 
      results, mutex).cost
    rightCost := CalculateNumericRange(
      LinearSearchRange{r.x + k + 1, r.n - k - 1}, results, mutex).cost
    cost = uint64(r.x + k)
    if leftCost > rightCost {
      cost += leftCost 
    } else {
      cost += rightCost
    } 

    switch { 
    case cost == minCost:
      minSplitPoints = append(minSplitPoints, k)
    case cost < minCost:
      minCost = cost
      minSplitPoints = []int{k}
    }
  }

  result := LinearSearchResult{minCost, minSplitPoints}
  mutex.Lock()
  (*results)[r] = result
  mutex.Unlock()

  return result


}

func CalculateNumericC(n int,
  results *map[LinearSearchRange]LinearSearchResult,
  mutex *sync.Mutex) LinearSearchResult {
  return CalculateNumericF(1, n-1, results, mutex)
}

func CalculateNumericF(x int, n int, 
  results *map[LinearSearchRange]LinearSearchResult,
  mutex *sync.Mutex) LinearSearchResult {
  return CalculateNumericRange(LinearSearchRange{x, n}, results, mutex)
}
