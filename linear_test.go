package searchcost

import "testing"
import "fmt"

var formatTests = []struct {
  val Linear
  str string
}{
  {Linear{0,0}, "0"},
  {Linear{1,0}, "x"},
  {Linear{0,1}, "1"},
  {Linear{1,1}, "x+1"},
  {Linear{2,0}, "2x"},
  {Linear{0,2}, "2"},
  {Linear{3,5}, "3x+5"},
}

func TestLinearString(t *testing.T) {
  for _, test := range formatTests {
    if test.val.String() != test.str {
      t.Error(fmt.Sprintf("Expected %s, was %s\n", 
        test.str, test.val.String()))
    }
  }
}

var evalTests = []struct {
  f     Linear
  x     uint64
  total uint64 // = ax + b
}{
  {Linear{0,0}, 0, 0},
  {Linear{0,0}, 6, 0},
  {Linear{0,7}, 0, 7},
  {Linear{0,7}, 6, 7},
  {Linear{3,0}, 0, 0},
  {Linear{3,0}, 4, 12},
  {Linear{4,5}, 0, 5},
  {Linear{4,5}, 4, 21},
}

func TestLinearEval(t *testing.T) { 
  for _, test := range evalTests {
    if test.f.Eval(test.x) != test.total {
      t.Error(fmt.Sprintf("f(x)=%s, f(%d)!=%d\n",
        test.f.String(), test.x, test.total));
    }
  }
}

var compareTests = []struct {
  a      Linear
  b      Linear 
  expect LinearCompare
}{
  {Linear{2,3}, Linear{2,3}, LINEAR_COMPARE_EQUAL},
  {Linear{0,3}, Linear{0,3}, LINEAR_COMPARE_EQUAL},
  {Linear{3,0}, Linear{3,0}, LINEAR_COMPARE_EQUAL},
  {Linear{0,2}, Linear{0,3}, LINEAR_COMPARE_LESS_OR_EQUAL},
  {Linear{5,2}, Linear{5,3}, LINEAR_COMPARE_LESS_OR_EQUAL},
  {Linear{4,6}, Linear{5,5}, LINEAR_COMPARE_LESS_OR_EQUAL},
  {Linear{5,3}, Linear{5,2}, LINEAR_COMPARE_GREATER_OR_EQUAL},
  {Linear{0,3}, Linear{0,2}, LINEAR_COMPARE_GREATER_OR_EQUAL},
  {Linear{4,5}, Linear{3,6}, LINEAR_COMPARE_GREATER_OR_EQUAL},
  {Linear{4,6}, Linear{3,6}, LINEAR_COMPARE_GREATER_OR_EQUAL},
  {Linear{4,4}, Linear{3,6}, LINEAR_COMPARE_INTERSECTS},
  {Linear{5,6}, Linear{3,9}, LINEAR_COMPARE_INTERSECTS},
  {Linear{3,9}, Linear{5,6}, LINEAR_COMPARE_INTERSECTS},
}

func TestLinearCompare(t *testing.T) {
  for _, test := range compareTests {
    cmp := test.a.Compare(test.b)
    if test.expect != cmp { 
      t.Error(fmt.Sprintf("%s vs %s expected %d, actual %d",
        test.a.String(), test.b.String(), test.expect, cmp))
    }
  }
}

var compareFromTests = []struct {
  a      Linear
  b      Linear
  n      uint64
  expect LinearCompare
}{
}

func TestLinearCompareFrom(t *testing.T) {
  for _, test := range compareFromTests {
    cmp := test.a.CompareFrom(test.b, test.n)
    if test.expect != cmp {
      t.Error(fmt.Sprintf("%s vs %s for x >= %d: expected %d, actual %d",
        test.a.String(), test.b.String(), test.n, test.expect, cmp))
    }
  }
}

var compareLinearIntersections = []struct {
  a Linear
  b Linear
  x uint64
}{
  {Linear{5,7}, Linear{3,19}, 6},
  {Linear{5,7}, Linear{3,19}, 6},
  {Linear{6,20}, Linear{4,14}, 1},
  {Linear{5,7}, Linear{3,18}, 5},
}

func TestLinearIntersections(t *testing.T) {
  for _, test := range compareLinearIntersections {
    tx := test.a.Intersection(test.b)
    if tx != test.x {
      t.Error(fmt.Sprintf("%s intersection %s, expect %d (was %d)",
        test.a, test.b, test.x, tx))
    }
  } 
}
