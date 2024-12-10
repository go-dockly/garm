package main

import (
	"fmt"
)

type geometry interface {
	area() float64
}

type rect struct {
	width, height float64
}

func (r rect) area() float64 {
	return r.width * r.height
}

func main() {
	r := rect{width: 3, height: 4}
	fmt.Println(r.area())
}
