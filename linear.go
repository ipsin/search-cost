package searchcost

import "fmt"

// Represents a linear value (ax+b) which is defined for integers x >= 1.
type Linear struct {
  a,b int64
}

type LinearCompare int

const (
  LINEAR_COMPARE_LESS_OR_EQUAL = iota
  LINEAR_COMPARE_EQUAL
  LINEAR_COMPARE_GREATER_OR_EQUAL
  LINEAR_COMPARE_INTERSECTS
)

func NewLinear(ta int64, tb int64) *Linear {
  return &Linear{ta,tb}
}

func (l *Linear) A() int64 {
  return l.a
}

func (l *Linear) B() int64 {
  return l.b
}

func (l *Linear) String() string {
  switch {
  case l.a == 0:
    return fmt.Sprintf("%d", l.b)
  case l.a == 1 && l.b == 0:
    return "x"
  case l.a == 1:
    if l.b < 0 {
      return fmt.Sprintf("x-%d", -l.b) 
    } else { 
      return fmt.Sprintf("x+%d", l.b) 
    }
  case l.b == 0:
    return fmt.Sprintf("%dx", l.a)
  }
  if l.b < 0 {
    return fmt.Sprintf("%dx%d", l.a, l.b)
  } else {
    return fmt.Sprintf("%dx+%d", l.a, l.b)
  }
}

func (l *Linear) Eval(x int64) int64 {
  return l.a * x + l.b
}

// Compare two linear values for all x >= 1
func (l *Linear) Compare(m *Linear) LinearCompare {
  return l.CompareFrom(m, 1)
}

func (l *Linear) Equal(m *Linear) bool {
  return l.a == m.a && l.b == m.b
}

// Compare two linear values for all x >= n
func (l *Linear) CompareFrom(m *Linear, n int64) LinearCompare {
  ln := l.Eval(n)
  mn := m.Eval(n)

  switch {
  case l.a == m.a && l.b == m.b:
    return LINEAR_COMPARE_EQUAL
  case ln >= mn && l.a >= m.a:
    return LINEAR_COMPARE_GREATER_OR_EQUAL
  case ln <= mn && l.a <= m.a:
    return LINEAR_COMPARE_LESS_OR_EQUAL
  }
  return LINEAR_COMPARE_INTERSECTS
}

// Compare two linear values for s <= x <= t.  Requires 1 <= s <= t.
func (l *Linear) CompareBetween(m *Linear, s int64, t int64) LinearCompare {
  ls := l.Eval(s)
  lt := l.Eval(t)
  ms := m.Eval(s)
  mt := m.Eval(t)

  switch {
  case ls == ms && lt == mt:
    return LINEAR_COMPARE_EQUAL
  case ls >= ms && lt >= mt:
    return LINEAR_COMPARE_GREATER_OR_EQUAL
  case ls <= ms && lt <= mt:
    return LINEAR_COMPARE_LESS_OR_EQUAL
  }
  return LINEAR_COMPARE_INTERSECTS
}

// Returns the intersection of two lines, rounded down.  This will crash
// if both have the same slope.  If the intersection point is x < 1, the
// value 1 will be returned.
func (l *Linear) Intersection(m *Linear) int64 {
  numPos := m.b >= l.b
  denPos := l.a > m.a
  var num, den int64

  if numPos { 
    num = m.b - l.b
  } else {
    num = l.b - m.b
  }

  if denPos { 
    den = l.a - m.a
  } else {
    den = m.a - l.a
  }

  xIntercept := num / den
  if numPos == denPos && xIntercept > 0 { 
    return xIntercept
  }

  return int64(1)
}
