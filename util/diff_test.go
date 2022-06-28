package util

import "testing"

func TestClear(t *testing.T) {
	a := `a
b
c
d
e
f`
	b := ``
	t.Log(Diff(a, b))
}

func TestInit(t *testing.T) {
	a := ``
	b := `a
b
c
d
e
f`
	t.Log(Diff(a, b))
}

func TestReplace(t *testing.T) {
	a := `a
b
c
d
e
f`
	b := `h
i
j
k
l
m`
	t.Log(Diff(a, b))
}

func TestDiff(t *testing.T) {
	a := `a
你好
c
d
e
f
g`
	b := `a
h
b
d
e
h
j
l`
	t.Log(Diff(a, b))
}