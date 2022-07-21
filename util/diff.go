package util

import (
	"bytes"
	"fmt"
	"strings"
)

func Diff(src, dest string) string {
	diff := NewLineDiffable(src, dest)
	es := Myers(diff)
	buf := bytes.NewBuffer(nil)
	x, y := 0, 0
	for i := range es {
		switch es[i] {
		case OpMov:
			x++
			y++
		case OpInsert:
			fmt.Fprintln(buf, "+", diff.Y(y), "<br />")
			y++
		case OpDel:
			fmt.Fprintln(buf, "-", diff.X(x), "<br />")
			x++
		}
	}
	return buf.String()
}

func DiffToArr(src, dest string) (arr []string) {
	diff := NewLineDiffable(src, dest)
	es := Myers(diff)
	x, y := 0, 0
	for i := range es {
		switch es[i] {
		case OpMov:
			arr = append(arr, "o "+diff.Y(y))
			x++
			y++
		case OpInsert:
			arr = append(arr, "+ "+diff.Y(y))
			y++
		case OpDel:
			arr = append(arr, "- "+diff.X(x))
			x++
		}
	}
	return
}

type EditScript []Op

func (e EditScript) reverse() EditScript {
	es := make(EditScript, 0, len(e))
	for i := len(e) - 1; i >= 0; i-- {
		es = append(es, e[i])
	}
	return es
}

type Op int8

const (
	OpInsert = 1
	OpMov    = 0
	OpDel    = -1
)

type LineDiffable struct {
	aText []string
	bText []string
}

func NewLineDiffable(a, b string) *LineDiffable {
	var aText, bText []string
	if a != "" {
		aText = strings.Split(a, "\n")
	}
	if b != "" {
		bText = strings.Split(b, "\n")
	}

	return &LineDiffable{
		aText: aText,
		bText: bText,
	}
}

func (l *LineDiffable) LenA() int {
	return len(l.aText)
}

func (l *LineDiffable) LenB() int {
	return len(l.bText)
}

func (l *LineDiffable) Equal(x, y int) bool {
	return l.aText[x] == l.bText[y]
}

func (l *LineDiffable) X(x int) string {
	return l.aText[x]
}

func (l *LineDiffable) Y(y int) string {
	return l.bText[y]
}

type Diffable interface {
	LenA() int           // length of src
	LenB() int           // length of dest
	Equal(int, int) bool // compare src[ai] and dest[bi]
}

func Myers(ab Diffable) EditScript {
	aLen := ab.LenA()
	bLen := ab.LenB()

	// preCheck
	if aLen == 0 {
		script := EditScript{}
		for i := 0; i < bLen; i++ {
			script = append(script, OpInsert)
		}
		return script
	} else if bLen == 0 {
		script := EditScript{}
		for i := 0; i < aLen; i++ {
			script = append(script, OpDel)
		}

		return script
	}

	max := aLen + bLen
	v := make([]int, 2*max+1)
	trace := make([][]int, 0)
search:
	// myers
	for d := 0; d <= max; d++ {
		vc := make([]int, 2*max+1)
		copy(vc, v)
		trace = append(trace, vc)

		for k := -d; k <= d; k += 2 {
			if k < -bLen || k > aLen {
				continue
			}
			var x int
			if k == -d || (k != d && v[max+k-1] < v[max+k+1]) {
				x = v[max+k+1]
			} else {
				x = v[max+k-1] + 1
			}

			y := x - k
			for x < aLen && y < bLen && ab.Equal(x, y) {
				x++
				y++
			}

			v[max+k] = x

			if x == aLen && y == bLen {
				break search
			}
		}
	}

	// 完全替换
	if len(trace)-1 == max {
		script := EditScript{}
		for i := 0; i < aLen; i++ {
			script = append(script, OpDel)
		}
		for i := 0; i < bLen; i++ {
			script = append(script, OpInsert)
		}
		return script
	}

	x := aLen
	y := bLen
	script := EditScript{}

	// 回溯
	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y

		var prevk int

		if k == -d || (k != d && v[max+k-1] < v[max+k+1]) {
			prevk = k + 1
		} else {
			prevk = k - 1
		}

		prevx := v[max+prevk]
		prevy := prevx - prevk
		for x > prevx && y > prevy {
			script = append(script, OpMov)
			x--
			y--
		}

		if d > 0 {
			if x == prevx {
				script = append(script, OpInsert)
			} else {
				script = append(script, OpDel)
			}
		}

		x, y = prevx, prevy
	}

	return script.reverse()
}
