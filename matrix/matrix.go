package matrix

import (
	"errors"
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

type Dimension struct {
	Width  int
	Height int
}

type Matrix[T constraints.Integer | constraints.Float] struct {
	Dimensions Dimension
	Values     [][]T
}

// New instantiates a Matrix with the passed values.
func New[T constraints.Integer | constraints.Float](err *error, values ...[]T) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	height := len(values)
	if height < 1 {
		*err = errors.New("cannot create a Matrix with a dimension that is less than 1")
		return Matrix[T]{}
	}

	width := len(values[0])
	if width < 1 {
		*err = errors.New("cannot create a Matrix with a dimension that is less than 1")
		return Matrix[T]{}
	}

	for _, row := range values {
		if len(row) != width {
			*err = errors.New("cannot create a Matrix with different row lengths")
			return Matrix[T]{}
		}
	}

	return Matrix[T]{
		Dimensions: Dimension{
			Width:  width,
			Height: height,
		},
		Values: values,
	}
}

// NewZero instantiates a zero matrix of the specified dimensions.
func NewZero[T constraints.Integer | constraints.Float](err *error, dim Dimension) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	if dim.Width < 1 || dim.Height < 1 {
		*err = errors.New("cannot create a Matrix with a dimension that is less than 1")
		return Matrix[T]{}
	}

	values := make([][]T, dim.Height)
	for y := 0; y < dim.Height; y++ {
		values[y] = make([]T, dim.Width)
	}

	return Matrix[T]{
		Dimensions: Dimension{
			Width:  dim.Width,
			Height: dim.Height,
		},
		Values: values,
	}
}

// NewIdentity instantiates an identity matrix of the specified dimensions.
func NewIdentity[T constraints.Integer | constraints.Float](err *error, dim Dimension) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	if dim.Width < 1 || dim.Height < 1 {
		*err = errors.New("cannot create a Matrix with a dimension that is less than 1")
		return Matrix[T]{}
	}

	// Check the matrix is square.
	if dim.Height != dim.Width {
		*err = errors.New("an identity matrix must be square")
		return Matrix[T]{}
	}

	values := make([][]T, dim.Height)
	for y := 0; y < dim.Height; y++ {
		values[y] = make([]T, dim.Width)
		values[y][y] = 1
	}

	return Matrix[T]{
		Dimensions: Dimension{
			Width:  dim.Width,
			Height: dim.Height,
		},
		Values: values,
	}
}

// MultiplyScalar multiplies the Matrix by a scalar.
func (a Matrix[T]) MultiplyScalar(err *error, x T) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	m := a.Clone()
	for i := 0; i < a.Dimensions.Width; i++ {
		for j := 0; j < a.Dimensions.Height; j++ {
			m.Values[j][i] *= x
		}
	}
	return m
}

// Multiply multiplies the matrix by another matrix.
// The height of matix B must match the width of matrix A.
func (a Matrix[T]) Multiply(err *error, b Matrix[T]) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	// Check matrices can be multiplied.
	if a.Dimensions.Width != b.Dimensions.Height {
		*err = errors.New("cannot multiply matrices due to incompatible dimensions")
		return Matrix[T]{}
	}

	// Initialize the resulting matrix.
	width := b.Dimensions.Width
	height := a.Dimensions.Height
	m := make([][]T, height)
	for j := 0; j < height; j++ {
		m[j] = make([]T, width)
	}

	// Iterate rows of matrix a.
	for j := 0; j < a.Dimensions.Height; j++ {
		// Iterate through coumns of matrix b.
		for x := 0; x < b.Dimensions.Width; x++ {
			// Calculate dot product of matrix cells.
			sum := T(0)
			for i := 0; i < a.Dimensions.Width; i++ {
				sum += (a.Values[j][i] * b.Values[i][x])
			}
			m[j][x] = sum
		}
	}

	return Matrix[T]{
		Dimensions: Dimension{
			Width:  width,
			Height: height,
		},
		Values: m,
	}
}

// Add a matrix to another one.
// The dimensions of the matrices must match.
func (a Matrix[T]) Add(err *error, b Matrix[T]) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	// Check matrices can be added.
	if a.Dimensions.Width != b.Dimensions.Width || a.Dimensions.Height != b.Dimensions.Height {
		*err = errors.New("cannot add matrices due to incompatible dimensions")
		return Matrix[T]{}
	}

	m := a.Clone()
	for j := 0; j < a.Dimensions.Height; j++ {
		for i := 0; i < a.Dimensions.Width; i++ {
			m.Values[j][i] += b.Values[j][i]
		}
	}

	return m
}

// Subtract a matrix from another one.
// The dimensions of the matrices must match.
func (a Matrix[T]) Subtract(err *error, b Matrix[T]) Matrix[T] {
	// Avoid hiding previous errors
	if *err != nil {
		return Matrix[T]{}
	}

	// Check matrices can be subtracted.
	if a.Dimensions.Width != b.Dimensions.Width || a.Dimensions.Height != b.Dimensions.Height {
		*err = errors.New("cannot subtract matrices due to incompatible dimensions")
		return Matrix[T]{}
	}

	m := a.Clone()
	for j := 0; j < a.Dimensions.Height; j++ {
		for i := 0; i < a.Dimensions.Width; i++ {
			m.Values[j][i] -= b.Values[j][i]
		}
	}

	return m
}

// Inverse a matrix
// TODO: figure out how to handle integers properly.
func (a Matrix[T]) Inverse(err *error) Matrix[T] {
	// Avoid hiding previous errors.
	if *err != nil {
		return Matrix[T]{}
	}

	// Check the matrix is square.
	if a.Dimensions.Height != a.Dimensions.Width {
		*err = errors.New("cannot calculate the inverse of a non-square matrix")
		return Matrix[T]{}
	}

	if a.Dimensions.Width == 2 {
		// Invert a 2x2 matrix.
		// https://en.wikipedia.org/wiki/Invertible_matrix#Inversion_of_2_%C3%97_2_matrices
		detA := (a.Values[0][0] * a.Values[1][1]) - (a.Values[0][1] * a.Values[1][0])
		if detA == 0 {
			*err = errors.New("cannot invert, matrix is singular")
			return Matrix[T]{}
		}
		detA = 1 / detA
		m := New(
			err,
			[]T{a.Values[1][1], -a.Values[0][1]},
			[]T{-a.Values[1][0], a.Values[0][0]},
		)
		if *err != nil {
			return Matrix[T]{}
		}
		return m.MultiplyScalar(err, detA)
	} else {
		// Invert a n-n matrix using Gausian Elimination.
		*err = errors.New("Inverse not implemented for matrixes width dimensions above 2x2")
		return Matrix[T]{}
	}
}

// Transpose calculates the transpose of a Matrix.
func (a Matrix[T]) Transpose(err *error) Matrix[T] {
	values := make([][]T, a.Dimensions.Width)
	for j := 0; j < a.Dimensions.Width; j++ {
		values[j] = make([]T, a.Dimensions.Height)
		for i := 0; i < a.Dimensions.Height; i++ {
			values[j][i] = a.Values[i][j]
		}
	}

	return Matrix[T]{
		Dimensions: Dimension{
			Width:  a.Dimensions.Height,
			Height: a.Dimensions.Width,
		},
		Values: values,
	}
}

func (a Matrix[T]) Clone() Matrix[T] {
	values := make([][]T, a.Dimensions.Height)
	for j := 0; j < a.Dimensions.Height; j++ {
		values[j] = make([]T, a.Dimensions.Width)
		for i := 0; i < a.Dimensions.Width; i++ {
			values[j][i] = a.Values[j][i]
		}
	}

	return Matrix[T]{
		Dimensions: Dimension{
			Width:  a.Dimensions.Width,
			Height: a.Dimensions.Height,
		},
		Values: values,
	}
}

func (a Matrix[T]) Equal(b Matrix[T]) bool {
	if a.Dimensions.Height != b.Dimensions.Height {
		return false
	}
	if a.Dimensions.Width != b.Dimensions.Width {
		return false
	}
	for i := 0; i < a.Dimensions.Width; i++ {
		for j := 0; j < a.Dimensions.Height; j++ {
			if a.Values[j][i] != b.Values[j][i] {
				return false
			}
		}
	}
	return true
}

func (a Matrix[T]) ApproxEqual(b Matrix[T], errorMargin T) bool {
	if a.Dimensions.Height != b.Dimensions.Height {
		return false
	}
	if a.Dimensions.Width != b.Dimensions.Width {
		return false
	}
	for i := 0; i < a.Dimensions.Width; i++ {
		for j := 0; j < a.Dimensions.Height; j++ {
			if math.Abs(float64(a.Values[j][i])-float64(b.Values[j][i])) > float64(errorMargin) {
				return false
			}
		}
	}
	return true
}

func (a Matrix[T]) String() string {
	return fmt.Sprint(a.Values)
}
