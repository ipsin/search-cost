package searchcost

import "fmt"
import "reflect"
import "strconv"
import "math/rand"
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
  lowerBound int64
  f          Linear
}

type PiecewiseSearchCost struct {
  // Array of F(x,i) where i is the index of this array.
  fi []Piecewise
}

var ZERO_PIECEWISE = Piecewise{
  []PiecewiseSegment {
    PiecewiseSegment{1, Linear{0,0}},
  },
}

func NewPiecewise(vals ...int64) *Piecewise {
  count := len(vals) / 3
  pieces := make([]PiecewiseSegment, count)
  for i := 0; i < count; i++ {
    pieces[i] = *NewPiecewiseSegment(vals[3*i], vals[3*i+1], vals[3*i+2])
  }
  return &Piecewise{pieces}
}

func NewPiecewiseSegment(bound int64, la int64, lb int64) *PiecewiseSegment { 
  return &PiecewiseSegment{bound, Linear{la, lb}}
}

// Find the PiecewiseSegment 
func (p *Piecewise) ActiveSegment(x int64) int {
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

func (p *Piecewise) LastLowerBound() int64 {
  return p.segments[len(p.segments) - 1].lowerBound
}

func (p *Piecewise) Eval(x int64) int64 {
  return p.segments[p.ActiveSegment(x)].f.Eval(x)
}

func (p *Piecewise) Equal(q *Piecewise) bool {
  return reflect.DeepEqual(p,q)
}

// If p=f(x), return a piecewise that takes the value q=f(x+n).  This
// result can have fewer segments, if n is greater than the lower bounds of
// some of p's segments.
func (p *Piecewise) OffsetX(n int64) Piecewise {
  result := Piecewise{[]PiecewiseSegment{}}
  var nextBound int64

  for i := 0; i < len(p.segments) - 1; i++ {  
    nextBound = p.segments[i + 1].lowerBound
    if n + 1 < nextBound {
      nextLinear := Linear{p.segments[i].f.a,
        p.segments[i].f.b + n *p.segments[i].f.a}

      if len(result.segments) == 0 {
        result.segments = append(result.segments, 
          PiecewiseSegment{1, nextLinear})
      } else {
        result.segments = append(result.segments, 
          PiecewiseSegment{p.segments[i].lowerBound - n, nextLinear})
      }
    }
  }

  // Now append the last segment...
  lastSegment := p.segments[len(p.segments) - 1]
  nextLinear := Linear{lastSegment.f.a, lastSegment.f.b + n * lastSegment.f.a}
  if lastSegment.lowerBound < n + 1 {
    result.segments = append(result.segments, PiecewiseSegment{1, nextLinear})
  } else { 
    result.segments = append(result.segments, 
      PiecewiseSegment{lastSegment.lowerBound - n, nextLinear})
  }

  return result
}

// If p=f(x), return a piecewise that takes the value q=f(x)+n.
func (p *Piecewise) OffsetY(n int64) Piecewise { 
  result := Piecewise{make([]PiecewiseSegment, len(p.segments))}

  for i := 0; i < len(p.segments); i++ {  
    result.segments[i] = PiecewiseSegment{p.segments[i].lowerBound,
        Linear{p.segments[i].f.a, p.segments[i].f.b + n}}
  }

  return result
}

func (a *Piecewise) Add(b *Piecewise) Piecewise { 
  result := Piecewise{[]PiecewiseSegment{}}

  aIndex, aEnd := 0, len(a.segments) - 1
  bIndex, bEnd := 0, len(b.segments) - 1
  lastIntersection, nextIntersection := int64(1), int64(1)
  done := false

  for !done {
    curAIndex := aIndex
    curBIndex := bIndex

    done = advanceIndexes(a, b, &aIndex, &bIndex, aEnd, bEnd, 
      &nextIntersection)

    result.segments = append(result.segments,
      PiecewiseSegment{lastIntersection, 
        Linear{a.segments[curAIndex].f.a + b.segments[curBIndex].f.a, 
               a.segments[curAIndex].f.b + b.segments[curBIndex].f.b}})

    lastIntersection = nextIntersection
  }

  return result
}

func (a *Piecewise) Subtract(b *Piecewise) Piecewise {
  result := Piecewise{[]PiecewiseSegment{}}

  aIndex, aEnd := 0, len(a.segments) - 1
  bIndex, bEnd := 0, len(b.segments) - 1
  lastIntersection, nextIntersection := int64(1), int64(1)
  done := false
  var prevLinear *Linear = nil

  for !done {
    curAIndex := aIndex
    curBIndex := bIndex

    done = advanceIndexes(a, b, &aIndex, &bIndex, aEnd, bEnd,
      &nextIntersection)

    nextLinear := Linear{a.segments[curAIndex].f.a - b.segments[curBIndex].f.a,
      a.segments[curAIndex].f.b - b.segments[curBIndex].f.b}

    if prevLinear == nil || !prevLinear.Equal(&nextLinear) {
      result.segments = append(result.segments,
        PiecewiseSegment{lastIntersection, nextLinear})
    }

    lastIntersection = nextIntersection
    prevLinear = &nextLinear
  }

  return result
}

func (p *Piecewise) String() string {
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

func RandomPiecewise(minSegments int, maxSegments int, minStep int64,
  maxStep int64, minA int64, maxA int64, minB int64, 
  maxB int64) Piecewise {
   
  segCount := minSegments + rand.Intn(maxSegments - minSegments)
  currentBound := int64(1)
  result := Piecewise{make([]PiecewiseSegment, segCount)}

  for i := 0; i < segCount; i++ {  
    a := minA + rand.Int63n(int64(maxA - minA))
    b := minB + rand.Int63n(int64(maxB - minB))
    result.segments[i].f = Linear{a,b}
    result.segments[i].lowerBound = currentBound
    currentBound += int64((minStep + 
      rand.Int63n(int64(maxStep - minStep))))
  }

  return result
}


// Return a Piecewise that (for all integers x >= 1) takes on the 
// lesser of p(x) and q(x).  
func (p *Piecewise) Min(q *Piecewise) Piecewise {
  return compose(p, q, minMaxCompose(p, q, true))
}

// Return a Piecewise that (for all integers x >= 1) takes on the 
// greater of p(x) and q(x).  
func (p *Piecewise) Max(q *Piecewise) Piecewise {
  return compose(p, q, minMaxCompose(p, q, false))
}

func CreatePiecewiseSearchCost() PiecewiseSearchCost {
  return PiecewiseSearchCost{[]Piecewise{
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1,Linear{0,0}},
    },},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1,Linear{1,0}},
    },},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1,Linear{1,1}},
    },},
    Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1,Linear{2,2}},
    },},
  }}
}

// Produce a Linear with the same slope as the last Piecewise, and greater
// than or equal to it at all points.
func (p *Piecewise) UpperBound() Linear {
  seg_end := len(p.segments) - 1
  upper_a := p.segments[seg_end].f.a
  upper_b := p.segments[seg_end].f.b
  var start,end int64
  var startEval,endEval int64

  for i := 0; i < seg_end - 1; i++ { 
    start = p.segments[i].lowerBound
    end = p.segments[i + 1].lowerBound - 1
  
    startEval = p.segments[i].f.Eval(start) 
    endEval = p.segments[i].f.Eval(end) 

    if upper_a * start + upper_b < startEval {
      upper_b = startEval - upper_a * start
    }
    if upper_a * end + upper_b < endEval {
      upper_b = endEval - upper_a * end
    }
  }

  return Linear{upper_a, upper_b}
}

// Produce a Linear with the same slope as the last Piecewise, and less 
// than or equal to it at all points.
func (p *Piecewise) LowerBound() Linear {
  seg_end := len(p.segments) - 1
  lower_a := p.segments[seg_end].f.a
  lower_b := p.segments[seg_end].f.b
  var start,end int64
  var startEval,endEval int64

  for i := 0; i < seg_end - 1; i++ {
    start = p.segments[i].lowerBound
    end = p.segments[i + 1].lowerBound - 1

    startEval = p.segments[i].f.Eval(start)
    endEval = p.segments[i].f.Eval(end)

    if lower_a * start + lower_b > startEval { 
      lower_b = startEval - lower_a * start
    }

    if lower_a * end + lower_b > endEval {
      lower_b = endEval - lower_a * end
    }
  }

  return Linear{lower_a, lower_b}
}


// Grow until p.fi[v] is defined
func (p *PiecewiseSearchCost) Grow(v int) {
  for len(p.fi) <= v {
    p.GrowOnce()
  }
}

func (p *PiecewiseSearchCost) Cost(n int) *Piecewise {
  return &((*p).fi[n])
}

func (p *PiecewiseSearchCost) GrowOnce() {
  var minPiecewise *Piecewise = nil
  minHits := []int{}
  n := len(p.fi)

  for k := 1; k < n; k++ {
    mid := Piecewise{[]PiecewiseSegment{
      PiecewiseSegment{1,Linear{1,int64(k)}},
    }}
    left := p.fi[k-1]
    right := p.fi[n-k-1].OffsetX(int64(k+1))

    leftRightMax := left.Max(&right)

    sum := mid.Add(&leftRightMax)
    if minPiecewise == nil {
      minPiecewise = &sum
      minHits = []int{k,}
    } else {
      t := sum.Min(minPiecewise)
      if t.Equal(minPiecewise) {
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
  // fmt.Printf("F(x,%d) = %s\n", n, minPiecewise.String())

  // lb := minPiecewise.LowerBound()
  // ub := minPiecewise.UpperBound()
  // fmt.Printf("%s <= F(x,%d) <= %s\n", lb.String(), n, 
  // ub.String())
  p.fi = append(p.fi, *minPiecewise) 
}


// Used for calculating Min() and Max().  Given two Piecewise (A and B), 
// the composition moving left-to-right (starting with x=1) begins with 
// A if startA is true, B otherwise.  At every x in switchIndex, the 
// selected Piecewise is alternated.  switchIndex should be strictly
// increasing, and the first value should be > 1 (since the value at x=1 
// is already determined using the boolean).
type composePiecewise struct {
  startA      bool
  switchIndex []int64
}

// Produce a Piecewise using the given composePiecewise formula (see above)
func compose(a *Piecewise, b *Piecewise, comp composePiecewise) Piecewise {
  // In the trivial case where we never switch, return the desired Piecewise
  if len(comp.switchIndex) == 0 {
    if comp.startA {
      return *a
    } else {
      return *b
    }
  }

  result := Piecewise{[]PiecewiseSegment{},}
  fromA := comp.startA
  var fromPiecewise *Piecewise

  if fromA {
    fromPiecewise = a 
  } else {
    fromPiecewise = b
  }

  lastSwitchPoint := int64(1)
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
      fromPiecewise = a 
    } else {
      fromPiecewise = b
    }

    lastSwitchPoint = value
  }

  // Write the segment past the last switch point (unless it matches the
  // previous function)
  // Move to the next segment, and append them all until then end.
  is_init := true
  for seg := fromPiecewise.ActiveSegment(lastSwitchPoint);
      seg < len(fromPiecewise.segments); seg++ {
    if lastLinear == nil || fromPiecewise.segments[seg].f != *lastLinear {
      if is_init { 
        result.segments = append(result.segments,
          PiecewiseSegment{lastSwitchPoint, fromPiecewise.segments[seg].f})
      } else { 
        result.segments = append(result.segments,
          PiecewiseSegment{fromPiecewise.segments[seg].lowerBound, 
            fromPiecewise.segments[seg].f})
      }
      lastLinear = &fromPiecewise.segments[seg].f;
    }
    is_init=false
  }

  return result
}

// For two Piecewise functions, and computing a Min() or Max(), return true
// if the first (a) should be used, and false if the second(b) should be 
// used.
func minMaxTakeFromFirst(a, b *Piecewise, isMin bool) bool {
  aStart, bStart := a.Eval(1), b.Eval(1)

  switch {
  case aStart < bStart:
    return isMin

  case bStart < aStart:
    return !isMin

  // They intersect at x=1, so take the one with the smaller slope, or
  // pick a if they coincide.
  default:
    switch {
      case a.segments[0].f.a < b.segments[0].f.a:
        return isMin
      case b.segments[0].f.a > a.segments[0].f.a:
        return !isMin
      default:
        return true
    }
  }
}

// Advance the segment whose lowerBound occurs first, moving left to right,
// or advance both if the lowerBounds coincide.
func advanceIndexes(a, b *Piecewise, aIndex, bIndex *int, 
                    aEnd, bEnd int, nextIntersection *int64) bool {
  switch {
  case *aIndex < aEnd && *bIndex < bEnd:
    switch {
    case a.segments[*aIndex + 1].lowerBound ==
         b.segments[*bIndex + 1].lowerBound:
      *nextIntersection = a.segments[*aIndex + 1].lowerBound
      *aIndex++
      *bIndex++
    case a.segments[*aIndex + 1].lowerBound <
         b.segments[*bIndex + 1].lowerBound:
      *nextIntersection = a.segments[*aIndex + 1].lowerBound
      *aIndex++
    default:
      *nextIntersection = b.segments[*bIndex + 1].lowerBound
      *bIndex++
    }

  case *aIndex == aEnd && *bIndex < bEnd:
    if *bIndex < bEnd { 
      *nextIntersection = b.segments[*bIndex + 1].lowerBound
    } else {
      *nextIntersection = 0
    }
    *bIndex++

  case *bIndex == bEnd && *aIndex < aEnd:
    if *aIndex < aEnd { 
      *nextIntersection = a.segments[*aIndex + 1].lowerBound
    } else {
      *nextIntersection = 0
    }
    *aIndex++

  default:
    *nextIntersection = 0
    return true
  }

  return false
}

func insertIntersection(fa, fb Linear, firstIntersection int64, 
                        nextIntersection int64, comp *composePiecewise,
                        takeFromA *bool, isMin bool) {
  var lineCompare LinearCompare

  // Compare fa to fb over the given segment
  if nextIntersection == 0 {
    lineCompare = fa.CompareFrom(&fb, firstIntersection)
  } else {
    lineCompare = fa.CompareBetween(&fb, firstIntersection, 
                                    nextIntersection - 1) 
  }
  // fmt.Printf("Comparing between %d and %d, linefrom %d\n", firstIntersection,
    // nextIntersection, lineCompare)

  aIsMin := (isMin == *takeFromA)

  switch {
  case aIsMin && lineCompare == LINEAR_COMPARE_GREATER_OR_EQUAL:
    fallthrough
  case !aIsMin && lineCompare == LINEAR_COMPARE_LESS_OR_EQUAL:
    *takeFromA = !*takeFromA
    comp.switchIndex = append(comp.switchIndex, firstIntersection)

  case aIsMin && lineCompare == LINEAR_COMPARE_INTERSECTS:
    fallthrough
  case !aIsMin && lineCompare == LINEAR_COMPARE_INTERSECTS:
    lineIntersect := fa.Intersection(&fb)

    lineCompare = fa.CompareBetween(&fb, firstIntersection, lineIntersect)

    switch { 
    case aIsMin && lineCompare == LINEAR_COMPARE_GREATER_OR_EQUAL:
      fallthrough
    case !aIsMin && lineCompare == LINEAR_COMPARE_LESS_OR_EQUAL:
      *takeFromA = !*takeFromA
      comp.switchIndex = append(comp.switchIndex, firstIntersection)
    }

    // Append the switch point past the intersection
    if nextIntersection == 0 || lineIntersect + 1 < nextIntersection { 
      *takeFromA = !*takeFromA
      comp.switchIndex = append(comp.switchIndex, lineIntersect + 1)
    }
  }
}

// Returns a composePiecewise that can be used to produce Min(a,b)
func minMaxCompose(a, b *Piecewise, isMin bool) composePiecewise {
  takeFromA := minMaxTakeFromFirst(a, b, isMin)
 
  result := composePiecewise{takeFromA, []int64{}}
 
  aIndex, aEnd := 0, len(a.segments) - 1
  bIndex, bEnd := 0, len(b.segments) - 1
  lastIntersection, nextIntersection := int64(1), int64(1)
  done := false

  for !done { 
    curAIndex := aIndex
    curBIndex := bIndex

    done = advanceIndexes(a, b, &aIndex, &bIndex, aEnd, bEnd, 
      &nextIntersection)

    insertIntersection(a.segments[curAIndex].f, b.segments[curBIndex].f,
                       lastIntersection, nextIntersection, &result,
                       &takeFromA, isMin)

    lastIntersection = nextIntersection
  }

  return result
}
