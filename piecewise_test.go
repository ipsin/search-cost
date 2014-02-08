package searchcost

import "fmt"
import "math/rand"
import "reflect"
import "strconv"
import "testing"

func TestPiecewiseActiveSegment(t *testing.T) {
  v := Piecewise{[]PiecewiseSegment{
    PiecewiseSegment{1, Linear{4,5}},
    PiecewiseSegment{5, Linear{3,11}},
    PiecewiseSegment{9, Linear{2,21}},
    PiecewiseSegment{12, Linear{1,34}},
  }}

  for x := uint64(1); x <= 15; x++ {
    switch {
    case x < 5:
      if v.ActiveSegment(x) != 0{
        t.Error(fmt.Sprintf("Expected segment 0 at x=%d", x))
      }
    case x < 9:
      if v.ActiveSegment(x) != 1 {
        t.Error(fmt.Sprintf("Expected segment 1 at x=%d", x))
      }
    case x < 12:
      if v.ActiveSegment(x) != 2 {
        t.Error(fmt.Sprintf("Expected segment 2 at x=%d", x))
      }
    default:
      if v.ActiveSegment(x) != 3 {
        t.Error(fmt.Sprintf("Expected segment 3 at x=%d", x))
      }
    }
  }
}

var piecewiseStringTests = []struct {
  p      Piecewise
  expect string
}{
   { Piecewise{[]PiecewiseSegment{
     PiecewiseSegment{1, Linear{4,5}},
     PiecewiseSegment{5, Linear{3,11}},
     PiecewiseSegment{9, Linear{2,21}},
     PiecewiseSegment{12, Linear{1,34}},
     }}, 
     "4x+5 (1<=x<5), 3x+11 (5<=x<9), 2x+21 (9<=x<12), x+34 (x>=12)",
   },
   { Piecewise{[]PiecewiseSegment{
     PiecewiseSegment{1, Linear{2,3}},
     }},
     "2x+3 (x>=1)",
   },
}


func TestPiecewiseString(t *testing.T) { 
  for _, test := range piecewiseStringTests {
    cmp := test.p.String()
    if test.expect != cmp {
      t.Error(fmt.Sprintf("Format mismatch (expected vs actual)\n" +
        "%s\n%s\n", test.expect, cmp))
    }
  }
}

type piecewisePair struct {
  a Piecewise
  b Piecewise
}

var piecewiseMinMaxTests = []piecewisePair {
  { Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{2,3}},
    }},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{2,3}},
    }},
  },
  { Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{4,0}},
      PiecewiseSegment{5, Linear{3,5}},
    }},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{3,5}},
      PiecewiseSegment{5, Linear{4,0}},
    }},
  },
  { Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{7,12}},
      PiecewiseSegment{3, Linear{4,19}},
      PiecewiseSegment{7, Linear{2,32}},
    }},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{0,52}},
    }},
  },
  // { Piecewise{[]PiecewiseSegment{
      // PiecewiseSegment{1, Linear{5,3}},
      // PiecewiseSegment{9, Linear{6,1}},
    // }},
    // Piecewise{[]PiecewiseSegment{
      // PiecewiseSegment{1, Linear{7,0}},
      // PiecewiseSegment{9, Linear{7,1}},
    // }},
  // },
}


const checkDistancePastBound = 5;

func TestPiecewiseMin(t *testing.T) { 
  DoTestPiecewiseMinMax(t, piecewiseMinMaxTests, "Min", true)
}

func TestPiecewiseMax(t *testing.T) { 
  DoTestPiecewiseMinMax(t, piecewiseMinMaxTests, "Max", false)
}


func DoTestPiecewiseMinMax(t *testing.T, pairs []piecewisePair, 
  minMaxStr string, isMin bool) {
  // Test the solution's correctness from x=1 until a few values
  // past the larger lowerBound.
  var val Piecewise

  for _, test := range pairs {
    if (isMin) { 
      val = test.a.Min(test.b)
    } else {
      val = test.a.Max(test.b)
    }

    boundA := test.a.LastBound()
    boundB := test.b.LastBound()
    var lastCheck uint64

    if boundA < boundB {
      lastCheck = boundB
    } else {
      lastCheck = boundA
    }
    lastCheck += checkDistancePastBound
    
    for x := uint64(1); x <= lastCheck; x++ {
      va := test.a.Eval(x)
      vb := test.b.Eval(x)
      vc := val.Eval(x)
 
      if (va <= vb) == isMin {
        if va != vc {
          t.Error(fmt.Sprintf("%s(%s ;; %s)=%s at x=%d, expected %d, was " + 
                  "%d [[compose=%s]]\n", minMaxStr, test.a, test.b, val, 
                  x, va, vc, minMaxCompose(&test.a, &test.b, isMin)))
        }
      } else {
        if vb != vc {
          t.Error(fmt.Sprintf("%s(%s ;; %s)=%s at x=%d, expected %d, was " + 
                  "%d [[compose=%s]]\n", minMaxStr, test.a, test.b, val, 
                  x, vb, vc, minMaxCompose(&test.a, &test.b, isMin)))
        } 
      }
    }
  }
}

func TestRandomMinMax(t *testing.T) {
  rand.NewSource(99)

  for i := 0; i < 10000; i++ { 
    f1 := RandomPiecewise(1, 10, uint64(1), uint64(10), uint64(0), uint64(8), uint64(0), uint64(8))
    f2 := RandomPiecewise(1, 10, uint64(1), uint64(10), uint64(0), uint64(8), uint64(0), uint64(8))
    // fmt.Printf("[[[ f1=%s, f2=%s ]]]\n", f1.String(), f2.String())

    DoTestPiecewiseMinMax(t, []piecewisePair{piecewisePair{f1, f2}}, 
      "Min", true)
    DoTestPiecewiseMinMax(t, []piecewisePair{piecewisePair{f1, f2}}, 
      "Max", false)
  }
}


var minMaxComposeTests = []struct {
  a         Piecewise
  b         Piecewise
  expectMin composePiecewise
  expectMax composePiecewise
}{
  { Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{4,7}},
     }},
    Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{6,6}},
       PiecewiseSegment{3, Linear{5,0}},
       PiecewiseSegment{12, Linear{4,3}},
     }},
    composePiecewise{true, []uint64{3,8,12}}, 
    composePiecewise{false, []uint64{3,8,12}}, 
  },
}

func TestMinMaxCompose(t *testing.T) {
  for _, test := range minMaxComposeTests {
    result := minMaxCompose(&test.a, &test.b, true)
    if !reflect.DeepEqual(result, test.expectMin) {
      t.Error(fmt.Sprintf("minMaxCompose(%s;%s;%s) expected %s, was %s",
              test.a, test.b, "(Min)", test.expectMin, result))
    }

    result = minMaxCompose(&test.a, &test.b, false)
    if !reflect.DeepEqual(result, test.expectMax) {
      t.Error(fmt.Sprintf("minMaxCompose(%s,%s,%s) expected %s, was %s",
              test.a, test.b, "(Max)", test.expectMin, result))
    }
  }
}

var piecewisecomposeTests = []struct {
  a       Piecewise
  b       Piecewise
  compose composePiecewise
  expect  Piecewise
}{
   { Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{4,0}},
       PiecewiseSegment{5, Linear{3,5}},
     }},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{3,5}},
       PiecewiseSegment{5, Linear{4,0}},
     }},
     composePiecewise{false, []uint64{5}},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{3,5}},
     }},
   },
   { Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{7,2}},
       PiecewiseSegment{7, Linear{3,12}},
     }},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{5,14}},
       PiecewiseSegment{12, Linear{6,0}},
     }},
     composePiecewise{true, []uint64{5,10}},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{7,2}},
       PiecewiseSegment{5, Linear{5,14}},
       PiecewiseSegment{10, Linear{3,12}},
     }},
   },
   { Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{7,2}},
       PiecewiseSegment{5, Linear{3,12}},
       PiecewiseSegment{10, Linear{9,15}},
     }},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{3,9}},
       PiecewiseSegment{5, Linear{8,2}},
       PiecewiseSegment{12, Linear{7,1}},
     }},
     composePiecewise{false, []uint64{4,8,14,20}},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{3,9}},
       PiecewiseSegment{4, Linear{7,2}},
       PiecewiseSegment{5, Linear{3,12}},
       PiecewiseSegment{8, Linear{8,2}},
       PiecewiseSegment{12, Linear{7,1}},
       PiecewiseSegment{14, Linear{9,15}},
       PiecewiseSegment{20, Linear{7,1}},
     }},
   },

  { Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{5,3}},
      PiecewiseSegment{9, Linear{6,1}},
    }},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{7,0}},
      PiecewiseSegment{9, Linear{7,1}},
    }},
    composePiecewise{false, []uint64{2}},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1, Linear{7,0}},
      PiecewiseSegment{2, Linear{5,3}},
      PiecewiseSegment{9, Linear{6,1}},
    }},
  },
}

func TestCompose(t *testing.T) {
  for _, test := range piecewisecomposeTests {
    result := compose(test.a, test.b, test.compose)
    if !reflect.DeepEqual(result, test.expect) {
      t.Error(fmt.Sprintf("Compose(%s,%s) with %s expected %s, was %s", 
              test.a, test.b, test.compose, test.expect, result))
    }
  }
}

func findNextValue(val *map[int]Piecewise, n int) {
  var minPiecewise *Piecewise = nil
  minHits := []int{}

  // HACK: FIXME
  for k := 1; k < n; k++ {
    mid := Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1,Linear{1,uint64(k)}},
    }}
    left := ((*val)[k-1])
    right := (*val)[n-k-1].OffsetX(uint64(k+1))

    sum := mid.Add(left.Max(right))
    // fmt.Printf("%d: [%s] + [%s] + [%s] ==> %s\n", k, mid, left, right, sum)

    if minPiecewise == nil {
      minPiecewise = &sum
      minHits = []int{k,}
    } else {
      t := sum.Min(*minPiecewise)
      if t.Equal(*minPiecewise) { 
        minHits = append(minHits, k)
      } else {
        minHits = []int{}
      }

      minPiecewise = &t
    }
  }

  isNormal := false
  minHitsStr := make([]string, len(minHits))
  for i, _ := range(minHitsStr) { 
    minHitsStr[i] = strconv.Itoa(minHits[i])
    if minHits[i] == n - 2 {
      isNormal = true
    }
  }

  if !isNormal { 
    fmt.Printf("%d is WEIRD\n", n)
  }
  fmt.Printf("F(x,%d) = %s\n", n, (*minPiecewise).String())
  (*val)[n] = *minPiecewise
}

func TestIterfunc(t *testing.T) {
  piecewise_map := make(map[int]Piecewise)
  piecewise_map[0] = Piecewise{[]PiecewiseSegment{
    PiecewiseSegment{1,Linear{0,0}},
  },}
  piecewise_map[1] = Piecewise{[]PiecewiseSegment{
    PiecewiseSegment{1,Linear{1,0}},
  },}
  piecewise_map[2] = Piecewise{[]PiecewiseSegment{
    PiecewiseSegment{1,Linear{1,1}},
  },}
  piecewise_map[3] = Piecewise{[]PiecewiseSegment{
    PiecewiseSegment{1,Linear{2,2}},
  },}

  for t := 4; t <= 50; t++ {  
    findNextValue(&piecewise_map, t)
  }
}
