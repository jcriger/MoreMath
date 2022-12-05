package matrix

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	var err error

	a := New(&err, []int{0, 1}, []int{2, 3})
	assert.NilError(t, err)
	assert.Equal(t, a.Dimensions.Width, 2)
	assert.Equal(t, a.Dimensions.Height, 2)
	i := 0
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			assert.Equal(t, a.Values[y][x], i)
			i++
		}
	}

	a = New(&err, []int{0, 1}, []int{2, 3}, []int{4, 5})
	assert.NilError(t, err)
	assert.Equal(t, a.Dimensions.Width, 2)
	assert.Equal(t, a.Dimensions.Height, 3)
	i = 0
	for y := 0; y < 3; y++ {
		for x := 0; x < 2; x++ {
			assert.Equal(t, a.Values[y][x], i)
			i++
		}
	}
}

func TestNewEmptyHeight(t *testing.T) {
	var err error

	_ = New[int](&err)
	assert.ErrorContains(t, err, "cannot create a Matrix with a dimension that is less than 1")
}

func TestNewEmptyWidth(t *testing.T) {
	var err error

	_ = New(&err, []int{})
	assert.ErrorContains(t, err, "cannot create a Matrix with a dimension that is less than 1")
}

func TestNewBadDimensnions(t *testing.T) {
	var err error

	_ = New(&err, []int{1}, []int{2, 3})
	assert.ErrorContains(t, err, "cannot create a Matrix with different row lengths")
}

func TestNewZero(t *testing.T) {
	var err error

	a := NewZero[int](&err, Dimension{
		Width:  2,
		Height: 2,
	})
	assert.NilError(t, err)
	assert.Equal(t, a.Dimensions.Width, 2)
	assert.Equal(t, a.Dimensions.Height, 2)
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			assert.Equal(t, a.Values[y][x], 0)
		}
	}

	a = NewZero[int](&err, Dimension{
		Width:  2,
		Height: 3,
	})
	assert.NilError(t, err)
	assert.Equal(t, a.Dimensions.Width, 2)
	assert.Equal(t, a.Dimensions.Height, 3)
	for y := 0; y < 3; y++ {
		for x := 0; x < 2; x++ {
			assert.Equal(t, a.Values[y][x], 0)
		}
	}
}

func TestNewIdentity(t *testing.T) {
	var err error

	// 2x2
	a := NewIdentity[int](&err, Dimension{
		Width:  2,
		Height: 2,
	})
	assert.NilError(t, err)
	r := New(&err, []int{1, 0}, []int{0, 1})
	assert.NilError(t, err)
	assert.Check(t, a.Equal(r))

	// 2x3 - invalid
	a = NewIdentity[int](&err, Dimension{
		Width:  2,
		Height: 3,
	})
	assert.ErrorContains(t, err, "must be square")
	err = nil

	// 4x4
	b := NewIdentity[float64](&err, Dimension{
		Width:  4,
		Height: 4,
	})
	assert.NilError(t, err)
	s := New(&err,
		[]float64{1, 0, 0, 0},
		[]float64{0, 1, 0, 0},
		[]float64{0, 0, 1, 0},
		[]float64{0, 0, 0, 1},
	)
	assert.NilError(t, err)
	assert.Check(t, b.Equal(s))
}

func TestMultiplyScalar(t *testing.T) {
	var err error

	a := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	m := a.MultiplyScalar(&err, 10)
	assert.NilError(t, err)

	r := New(&err, []int{10, 20}, []int{30, 40}, []int{50, 60})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))
}

func TestMultiply(t *testing.T) {
	var err error

	// 2x2 * 2x3
	a := New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	b := New(&err, []int{1, 2, 3}, []int{4, 5, 6})
	assert.NilError(t, err)
	m := a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r := New(&err, []int{9, 12, 15}, []int{19, 26, 33})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 2x2 * 2x2
	a = New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{7, 10}, []int{15, 22})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x2 * 2x1
	a = New(&err, []int{1, 2})
	assert.NilError(t, err)
	b = New(&err, []int{1}, []int{2})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{5})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 2x1 * 1x2
	a = New(&err, []int{1}, []int{2})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{1, 2}, []int{2, 4})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x3 * 3x3
	a = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8, 9})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{30, 36, 42})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 2x3 * 3x1
	a = New(&err, []int{1, 2, 3}, []int{4, 5, 6})
	assert.NilError(t, err)
	b = New(&err, []int{1}, []int{2}, []int{3})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{14}, []int{32})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x1 * 1x3
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x1 * 1x1
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{1})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))
}

func TestAdd(t *testing.T) {
	var err error

	// 2x2 + 2x2
	a := New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	b := New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	m := a.Add(&err, b)
	assert.Check(t, err == nil)
	r := New(&err, []int{2, 4}, []int{6, 8})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x3 + 1x3
	a = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	m = a.Add(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{2, 4, 6})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 3x1 + 3x1
	a = New(&err, []int{1}, []int{2}, []int{3})
	assert.NilError(t, err)
	b = New(&err, []int{1}, []int{2}, []int{3})
	assert.NilError(t, err)
	m = a.Add(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{2}, []int{4}, []int{6})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x1 + 1x1
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1})
	assert.NilError(t, err)
	m = a.Add(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{2})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))
}

func TestSubtract(t *testing.T) {
	var err error

	// 2x2 - 2x2
	a := New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	b := New(&err, []int{4, 3}, []int{2, 1})
	assert.NilError(t, err)
	m := a.Subtract(&err, b)
	assert.Check(t, err == nil)
	r := New(&err, []int{-3, -1}, []int{1, 3})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x3 - 1x3
	a = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2, 3})
	assert.NilError(t, err)
	m = a.Subtract(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{0, 0, 0})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 3x1 - 3x1
	a = New(&err, []int{1}, []int{2}, []int{3})
	assert.NilError(t, err)
	b = New(&err, []int{3}, []int{2}, []int{1})
	assert.NilError(t, err)
	m = a.Subtract(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{-2}, []int{0}, []int{2})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	// 1x1 - 1x1
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1})
	assert.NilError(t, err)
	m = a.Subtract(&err, b)
	assert.Check(t, err == nil)
	r = New(&err, []int{0})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))
}

func TestInverse(t *testing.T) {
	var err error

	// 2x2
	a := New(&err, []float64{1, -1}, []float64{0, 2})
	assert.NilError(t, err)
	inv := a.Inverse(&err)
	assert.Check(t, err == nil)
	r := New(&err, []float64{1, 1.0 / 2.0}, []float64{0, 1.0 / 2.0})
	assert.NilError(t, err)
	assert.Check(t, inv.ApproxEqual(r, 0.0001))

	// 2x2 - validate A-1 * I = I
	r = inv.Multiply(&err, a)
	assert.NilError(t, err)
	i := NewIdentity[float64](&err, a.Dimensions)
	assert.NilError(t, err)
	assert.Check(t, r.ApproxEqual(i, 0.0001))

	// 2x2
	a = New(&err, []float64{1, 2}, []float64{3, -5})
	assert.NilError(t, err)
	inv = a.Inverse(&err)
	assert.Check(t, err == nil)
	r = New(&err, []float64{5.0 / 11.0, 2.0 / 11.0}, []float64{3.0 / 11.0, -1.0 / 11.0})
	assert.NilError(t, err)
	assert.Check(t, inv.ApproxEqual(r, 0.0001))

	// 2x2 - validate A-1 * I = I
	r = inv.Multiply(&err, a)
	assert.NilError(t, err)
	i = NewIdentity[float64](&err, a.Dimensions)
	assert.NilError(t, err)
	assert.Check(t, r.ApproxEqual(i, 0.0001))
}

func TestOperationErrors(t *testing.T) {
	// TODO: add additional functions
	// TODO: improve testing of function chaining
	var err error

	// 2x2 * 3x1
	a := New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	b := New(&err, []int{1}, []int{2}, []int{3})
	assert.NilError(t, err)
	m := a.Multiply(&err, b)
	assert.ErrorContains(t, err, "cannot multiply matrices due to incompatible dimensions")
	assert.Check(t, m.Equal(Matrix[int]{}))

	m = a.Add(&err, b)
	assert.ErrorContains(t, err, "cannot multiply matrices due to incompatible dimensions")
	assert.Check(t, m.Equal(Matrix[int]{}))

	err = nil
	// 1x1 * 2x2
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	m = a.Multiply(&err, b)
	assert.ErrorContains(t, err, "cannot multiply matrices due to incompatible dimensions")
	assert.Check(t, m.Equal(Matrix[int]{}))

	err = nil
	// 1x1 + 2x2
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	m = a.Add(&err, b)
	assert.ErrorContains(t, err, "cannot add matrices due to incompatible dimensions")
	assert.Check(t, m.Equal(Matrix[int]{}))

	err = nil
	// 1x1 - 2x2
	a = New(&err, []int{1})
	assert.NilError(t, err)
	b = New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	m = a.Subtract(&err, b)
	assert.ErrorContains(t, err, "cannot subtract matrices due to incompatible dimensions")
	assert.Check(t, m.Equal(Matrix[int]{}))
}

func TestTranspose(t *testing.T) {
	var err error

	a := New(&err, []int{1, 2})
	assert.NilError(t, err)
	m := a.Transpose(&err)

	r := New(&err, []int{1}, []int{2})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	a = New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	m = a.Transpose(&err)

	r = New(&err, []int{1, 3}, []int{2, 4})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))

	a = New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	m = a.Transpose(&err)

	r = New(&err, []int{1, 3, 5}, []int{2, 4, 6})
	assert.NilError(t, err)
	assert.Check(t, m.Equal(r))
}

func TestClone(t *testing.T) {
	var err error

	a := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)

	m := a.Clone()
	r := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)

	a.Dimensions.Width = 10
	assert.Check(t, m.Equal(r))

	a.Dimensions.Height = 10
	assert.Check(t, m.Equal(r))

	a.Values[0][0] = 10
	assert.Check(t, m.Equal(r))
}

func TestEqual(t *testing.T) {
	var err error

	a := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	assert.Check(t, a.Equal(a))

	b := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	assert.Check(t, a.Equal(b))

	c := New(&err, []int{2, 1}, []int{4, 3}, []int{6, 5})
	assert.NilError(t, err)
	assert.Check(t, !a.Equal(c))

	d := New(&err, []int{1, 2}, []int{3, 4})
	assert.NilError(t, err)
	assert.Check(t, !a.Equal(d))

	e := New(&err, []int{1, 2, 3}, []int{3, 4, 5}, []int{5, 6, 7})
	assert.NilError(t, err)
	assert.Check(t, !a.Equal(e))
}

func TestApproxEqual(t *testing.T) {
	var err error

	a := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	assert.Check(t, a.ApproxEqual(a, 0))

	b := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	assert.Check(t, a.ApproxEqual(b, 0))

	b = New(&err, []int{2, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	assert.Check(t, !a.ApproxEqual(b, 0))

	c := New(&err, []float64{2, 1}, []float64{4, 3}, []float64{6, 5})
	assert.NilError(t, err)
	d := New(&err, []float64{2, 1}, []float64{4, 3}, []float64{6, 5})
	assert.NilError(t, err)
	assert.Check(t, c.ApproxEqual(d, 0.0001))

	c = New(&err, []float64{1.0 / 2.0, 1.0 / 11.0}, []float64{2.0 / 3.55, 1.0 / 3.0})
	assert.NilError(t, err)
	d = New(&err, []float64{1.0 / 2.0, 1.0 / 11.0}, []float64{2.0 / 3.55, 1.0 / 3.0})
	assert.NilError(t, err)
	assert.Check(t, c.ApproxEqual(d, 0.0001))

	c = New(&err, []float64{1.0 / 2.0, 1.0 / 11.0})
	assert.NilError(t, err)
	d = New(&err, []float64{1.0 / 2.1, 1.1 / 11.0})
	assert.NilError(t, err)
	assert.Check(t, !c.ApproxEqual(d, 0.0001))
}

func TestString(t *testing.T) {
	var err error

	m := New(&err, []int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.NilError(t, err)
	assert.Equal(t, m.String(), "[[1 2] [3 4] [5 6]]")
}
