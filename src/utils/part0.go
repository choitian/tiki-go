package utils

import (
	"math"
)

var fib map[int]int = make(map[int]int)

// fibonacci is a function that returns
// a function that returns an int.

func Fibonacci(x int, fast bool) int {
	fast = fast
	if x == 0 {
		return 0
	} else if x == 1 {
		return 1
	} else if x == 2 {
		return 1
	} else {
		if fast {
			v, exist := fib[x]
			if exist {
				return v
			} else {
				ret := Fibonacci(x-1, fast) + Fibonacci(x-2, fast)
				fib[x] = ret
				return ret
			}
		} else {
			return Fibonacci(x-1, fast) + Fibonacci(x-2, fast)
		}
	}
}

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func (v Vertex) Scale2(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

type I interface {
	M()
}
