package searchcost

import "fmt"
import "reflect"
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

var piecewiseMinMaxTests = []struct {
  a           Piecewise
  b           Piecewise
}{
   // { Piecewise{[]PiecewiseSegment{
       // PiecewiseSegment{1, Linear{2,3}},
     // }},
     // Piecewise{[]PiecewiseSegment{
       // PiecewiseSegment{1, Linear{2,3}},
     // }},
   // },
   // { Piecewise{[]PiecewiseSegment{
       // PiecewiseSegment{1, Linear{4,0}},
       // PiecewiseSegment{5, Linear{3,5}},
     // }},
     // Piecewise{[]PiecewiseSegment{
       // PiecewiseSegment{1, Linear{3,5}},
       // PiecewiseSegment{5, Linear{4,0}},
     // }},
   // },
   { Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{7,12}},
       PiecewiseSegment{3, Linear{4,19}},
       PiecewiseSegment{7, Linear{2,32}},
     }},
     Piecewise{[]PiecewiseSegment{
       PiecewiseSegment{1, Linear{0,52}},
     }},
   },


}

const checkDistancePastBound = 5;

func TestPiecewiseMin(t *testing.T) { 
  // DoTestPiecewiseMinMax(t, "Min", true)
}

func TestPiecewiseMax(t *testing.T) { 
  // DoTestPiecewiseMinMax(t, "Max", false)
}


func DoTestPiecewiseMinMax(t *testing.T, minMaxStr string, isMin bool) {
  // Test the solution's correctness from x=1 until a few values
  // past the larger lowerBound.
  for _, test := range piecewiseMinMaxTests {
    min := test.a.Min(test.b)

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
      vc := min.Eval(x)
 
      if (va <= vb) == isMin {
        if va != vc {
          t.Error(fmt.Sprintf("%s(%s ;; %s) at x=%d, expected %d, was %d\n", 
                  minMaxStr, test.a, test.b, x, va, vc))
        }
      } else {
        if vb != vc {
          t.Error(fmt.Sprintf("%s(%s ;; %s) at x=%d, expected %d, was %d\n", 
                  minMaxStr, test.a, test.b, x, vb, vc))
        } 
      }
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

