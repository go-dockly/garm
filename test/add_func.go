package main

func add(a, b int) int {
	var c = a + b
	return c
}

func main() {
	var a = 1
	var b = 2
	var c = add(a, b)
	_ = c
}
