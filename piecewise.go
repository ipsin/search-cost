package searchcost

import "fmt"
// import "math"
import "strings"

// A Piecewise is a list of linear functions (Linear), ordered by the 
// lower bound where that Linear takes effect.  The value of the Piecewise
// at n is the value of the Linear with the largest lowerBound less than
// or equal to n.  A Piecewise is defined for all integer x >= 1,
type Piecewise struct {
  segments []PiecewiseSegment
}

// Starting at lowerBound, this segment is equal to f.  This is true until
// the lowerBound of the next Piecewise (or for all x >= lowerBound, if it's
// the last segment.
type PiecewiseSegment struct {
  lowerBound uint64
  f          Linear
}

// Find the PiecewiseSegment 
func (p Piecewise) ActiveSegment(x uint64) int {
  // x = 1 will always be the start of the first segment.
  if x == 1 {
    return 0
  }

  // It will be common for x to surpass the highest lowerBound, so we 
  // optimize for this case
  high := len(p.segments) - 1
  if x >= p.segments[high].lowerBound {
    return high
  }


  low := 0
  mid := 0
  
  for low <= high { 
    mid = (low + high) / 2
    switch {
    case p.segments[mid].lowerBound < x:
      low = mid + 1
    case p.segments[mid].lowerBound > x:
      high = mid - 1
    default:
      return mid
    }
  }

  if high < low {
    tmp := high
    low = high
    high = tmp
  }

  if low >= 0 && p.segments[low].lowerBound < x {
    return low
  }
  if high >= 0 && p.segments[high].lowerBound < x {
    return low
  }

  return -1 
}

func (p Piecewise) LastBound() uint64 {
  return p.segments[len(p.segments) - 1].lowerBound
}

func (p Piecewise) Eval(x uint64) uint64 {
  return p.segments[p.ActiveSegment(x)].f.Eval(x)
}

func (p Piecewise) String() string {
  if p.segments == nil || len(p.segments) == 0 {
    return "[EMPTY PIECEWISE]"
  }
 

  lastSegment := len(p.segments) - 1
  strs := make([]string, len(p.segments))
  
  for i := 0; i < lastSegment; i++ {
    strs[i] = fmt.Sprintf("%s (%d<=x<%d)", p.segments[i].f.String(),
      p.segments[i].lowerBound, p.segments[i+1].lowerBound)
  }

  strs[lastSegment] = fmt.Sprintf("%s (x>=%d)", 
      p.segments[lastSegment].f.String(), 
      p.segments[lastSegment].lowerBound)

  return strings.Join(strs, ", ")
}


// Return a Piecewise that (for all integers x >= 1) takes on the 
// lesser of p(x) and q(x).  
func (p Piecewise) Min(q Piecewise) Piecewise {
  return p
}

// Return a Piecewise that (for all integers x >= 1) takes on the 
// greater of p(x) and q(x).  
func (p Piecewise) Max(q Piecewise) Piecewise {
  return p
}

// Used for calculating Min() and Max().  Given two Piecewise (A and B), 
// the composition moving left-to-right (starting with x=1) begins with 
// A if startA is true, B otherwise.  At every x in switchIndex, the 
// selected Piecewise is alternated.  switchIndex should be strictly
// increasing, and the first value should be > 1 (since the value at x=1 
// is already determined using the boolean).
type composePiecewise struct {
  startA      bool
  switchIndex []uint64
}

// Produce a Piecewise using the given composePiecewise formula (see above)
func compose(a Piecewise, b Piecewise, comp composePiecewise) Piecewise {
  // In the trivial case where we never switch, return the desired Piecewise
  if len(comp.switchIndex) == 0 {
    if comp.startA {
      return a
    } else {
      return b
    }
  }

  result := Piecewise{[]PiecewiseSegment{},}
  fromA := comp.startA
  var fromPiecewise *Piecewise

  if fromA {
    fromPiecewise = &a 
  } else {
    fromPiecewise = &b
  }

  lastSwitchPoint := uint64(1)
  var lastLinear *Linear = nil

  for _, value := range comp.switchIndex {
    firstSeg := true

    // Take all the PiecewiseSegment from the previous point to the
    // next switch point and write them to the result
    for seg := fromPiecewise.ActiveSegment(lastSwitchPoint); 
        seg < len(fromPiecewise.segments) && 
        fromPiecewise.segments[seg].lowerBound < value; seg++ {

      // If this is the first segment since the switch, start it at the
      // switch point, instead of its lower bound.  
      if firstSeg {
        if lastLinear == nil || fromPiecewise.segments[seg].f != *lastLinear {
          result.segments = append(result.segments, 
            PiecewiseSegment{lastSwitchPoint, fromPiecewise.segments[seg].f})
        }
        firstSeg = false
      } else {
        if lastLinear == nil || fromPiecewise.segments[seg].f != *lastLinear {
          result.segments = append(result.segments, 
            fromPiecewise.segments[seg])
        }
      }
 
      lastLinear = &fromPiecewise.segments[seg].f
      lastSwitchPoint = value
    }

    // After processing the segments up until the switch, swap the source
    // of new segments and continue reading
    fromA = !fromA
    if fromA { 
      fromPiecewise = &a 
    } else {
      fromPiecewise = &b
    }

    lastSwitchPoint = value
  }

  // Write the segment past the last switch point (unless it matches the
  // previous function)
  seg := fromPiecewise.ActiveSegment(lastSwitchPoint)
  if lastLinear == nil || fromPiecewise.segments[seg].f != *lastLinear {
    result.segments = append(result.segments,
      PiecewiseSegment{lastSwitchPoint, fromPiecewise.segments[seg].f})
  } 

  return result
}

// Returns a composePiecewise that can be used to produce Min(a,b)
func minCompose(a, b Piecewise) composePiecewise {
  var takeFrom *Piecewise
  var takeFromA bool

  aStart, bStart := a.Eval(1), b.Eval(1)

  switch {
  case aStart < bStart:
    takeFromA = true

  case bStart < aStart:
    takeFromA = false

  // They intersect at x=1, so take the one with the smaller slope, or 
  // pick a if they coincide.
  default:
    switch { 
      case a.segments[0].f.a < b.segments[0].f.a:
        takeFromA = true
      case b.segments[0].f.a > a.segments[0].f.a:
        takeFromA = false
      default:
        takeFromA = true
    }
  } 

  if takeFromA {
    takeFrom = &a
  } else {
    takeFrom = &b
  }
}

// Returns a composePiecewise that can be used to produce Max(a,b)
func maxCompose(a, b Piecewise) composePiecewise {
  return composePiecewise{true, []uint64{1}}
}
