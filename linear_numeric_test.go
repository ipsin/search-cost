package searchcost

import "fmt"
import "testing"
import "sync"


// C(n) = F(1,n-1) = the min cost of searching 1...n
// sum(C(n), n=1..100) = 17575
func TestProjectEulerSum(t *testing.T) {
  results := make(map[LinearSearchRange]LinearSearchResult)
  mutex := sync.Mutex{}
  const search_end = 100
  const expected_sum = uint64(17575)
  
  sum := uint64(0)
  for n := 1; n <= search_end; n++ {
    lsr := CalculateNumericC(n, &results, &mutex)
    sum += lsr.cost
  }
  if sum != expected_sum { 
    t.Error(fmt.Sprintf("Incorrect sum(C(i), 1, %d), expected %d, was %d",
      search_end, expected_sum, sum))
  }
}

var projectEulerValues = []struct {
  n       int
  expect  uint64
}{
  {1, 0},
  {2, 1},
  {3, 2},
  {8, 12},
  {100, 400},
}

// Tests the calculation of C(n) against known values, using a common cache.
func TestEulerValuesPrecompute(t *testing.T) { 
  results := make(map[LinearSearchRange]LinearSearchResult)
  mutex := sync.Mutex{}

  for _, test := range projectEulerValues {
    val := CalculateNumericC(test.n, &results, &mutex) 

    if test.expect != val.cost {
      t.Error(fmt.Sprintf("C(%d), expected %d, was %d", test.n, test.expect,
        val.cost))
    }
  }
}

// As above, but the cache is not shared between calls to compute C(n).
func TestEulerValuesDirect(t *testing.T) {
  for _, test := range projectEulerValues {
    results := make(map[LinearSearchRange]LinearSearchResult)
    mutex := sync.Mutex{}
    val := CalculateNumericC(test.n, &results, &mutex)

    if test.expect != val.cost {
      t.Error(fmt.Sprintf("C(%d), expected %d, was %d", test.n, test.expect,
        val.cost))
    }
  }
}
