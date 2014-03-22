package main

import "fmt"
import "os"
import "searchcost"
// import "github.com/ipsin/searchcost"

func write_gnuplot(filename string, funcname string, n int, tp *searchcost.Piecewise, max_val int64) {

  file, err := os.Create(filename)
  if err != nil { panic(err) }
  defer func() {
    if err := file.Close(); err != nil {
      panic(err)
    }
  }()

  max_val /= 25 
  if max_val == 0 { 
    max_val = 1
  }

  _, err = file.WriteString(fmt.Sprintf("set title \"%s\"\n", funcname))
  if err != nil { panic(err) }
  _, err = file.WriteString(fmt.Sprintf("set out \"%05d.png\"\n", n))
  if err != nil { panic(err) }
  _, err = file.WriteString("set terminal png font \"arial\" 30\n")
  if err != nil { panic(err) }
  _, err = file.WriteString("set terminal png size 1280,800\n")
  if err != nil { panic(err) }
  _, err = file.WriteString(fmt.Sprintf("set xrange [%d:%d]\n", 1, tp.LastLowerBound() + 10))
  if err != nil { panic(err) }
  _, err = file.WriteString("set style fill solid border -1\n")
  if err != nil { panic(err) }
  _, err = file.WriteString(fmt.Sprintf("set ytics %d\n", max_val))
  if err != nil { panic(err) }
  _, err = file.WriteString("unset key\n")
  if err != nil { panic(err) }
  _, err = file.WriteString(fmt.Sprintf("plot \"%05d.data\" with boxes\n", n))
  if err != nil { panic(err) }
  _, err = file.WriteString(fmt.Sprintf("!convert \"%05d.png\" \"%05d.gif\"\n", n, n))
  if err != nil { panic(err) }
  _, err = file.WriteString(fmt.Sprintf("!rm -f \"%05d.png\"\n", n))
  if err != nil { panic(err) }
}

func write_piecewise(name string, n int, tp *searchcost.Piecewise) int64 {
  file, err := os.Create(name)
  if err != nil { panic(err) }
  defer func() { 
    if err := file.Close(); err != nil {
      panic(err)
    }
  }()

  mv := int64(0)

  lastValue := tp.LastLowerBound() + 15
  for k := int64(1); k < lastValue; k++ { 
    tr := tp.Eval(k)
    if tr > mv {
      mv = tr
    }
    _, err := file.WriteString(fmt.Sprintf("%d %d\n", k, tr))
    if err != nil { panic(err) }
  }
  return mv
}

func lameZeroPiecewise(tp *searchcost.Piecewise) bool {
  lastValue := tp.LastLowerBound() + 2
  for k := int64(1); k < lastValue; k++ { 
    if tp.Eval(k) != 0 {
      return false
    }
  }
  return true
}


func main() {
  costs := searchcost.CreatePiecewiseSearchCost()

  for t := 1; t < 1000 ; t++ {
    costs.Grow(t+1)
    r1 := costs.Cost(t).OffsetX(1)
    diff := costs.Cost(t + 1).Subtract(&r1)
    if !diff.Equal(&searchcost.ZERO_PIECEWISE) && !lameZeroPiecewise(&diff) {
      data_filename := fmt.Sprintf("webz/gnuplot/%05d.data", t)
      max_val := write_piecewise(data_filename, t, &diff)
      gnu_filename := fmt.Sprintf("webz/gnuplot/%05d.gnuplot", t)
      func_string := fmt.Sprintf("F(x,%d)-F(x+1,%d)", t+1, t)
      write_gnuplot(gnu_filename, func_string, t, &diff, max_val)

      fmt.Printf("** F(x,%d)-F(x+1,%d)=%s\n", t+1, t, diff.String())
    }
  }
}
