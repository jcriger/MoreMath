package main

import (
	"fmt"
	"log"

	"github.com/jcriger/MoreMath/matrix"
)

func main() {
	var err error

	A := matrix.New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	B := matrix.New(&err, []int{1}, []int{2})

	// C = (2AB)'
	C := A.MultiplyScalar(&err, 2).Multiply(&err, B).Transpose(&err)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("(2AB)' =", C)

	D := matrix.New(&err, []float64{1.1, 2}, []float64{3, 4.1}, []float64{5.1, 6})
	E := matrix.New(&err, []float64{1.1}, []float64{2})

	// F = -2(DE)
	F := D.Multiply(&err, E).MultiplyScalar(&err, -2)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("-2DE =", F)
}
