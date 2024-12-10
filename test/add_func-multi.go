package test

func main() {
	a := 1         // Short lived, goes directly to x0
	b := 2         // Short lived, goes directly to x1
	c := add(a, b) // Result needed later, store in x19

	d := 3         // Short lived, goes directly to x0
	e := 4         // Short lived, goes directly to x1
	f := add(d, e) // Result needed later, store in x20

	g := add(c, f) // Use x19,x20 as inputs via x0,x1
}

func add(a, b int) int {
	var c = a + b
	return c
}
