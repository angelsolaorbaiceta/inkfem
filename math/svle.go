package math

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkmath/nums"
)

// SingleVarLinEq represents a single variable linear equation of the form:
//
//	y = a*x + b
type SingleVarLinEq struct {
	a, b         float64
	isHorizontal bool
}

func (svle SingleVarLinEq) YIntercept() float64 {
	return svle.b
}

func (svle SingleVarLinEq) Slope() float64 {
	return svle.a
}

// IsHorizontal returns true if the single variable linear equation is a
// horizontal line, that is, the slope is zero.
func (svle SingleVarLinEq) IsHorizontal() bool {
	return svle.isHorizontal
}

// MakeSVLE creates a new single variable linear equation from the given coefficients.
func MakeSVLE(a, b float64) SingleVarLinEq {
	isHorizontal := nums.FloatsEqual(a, 0.0)
	return SingleVarLinEq{a, b, isHorizontal}
}

// MakeSVLEFromPoints creates a new single value linear equation from two points.
// The two points must have different x values, otherwise an error is returned.
func MakeSVLEFromPoints(x1, y1, x2, y2 float64) (*SingleVarLinEq, error) {
	if nums.FloatsEqual(x1, x2) {
		return nil, fmt.Errorf("x1 and x2 must be different")
	}

	var (
		a            = (y2 - y1) / (x2 - x1)
		b            = y1 - a*x1
		isHorizontal = nums.FloatsEqual(a, 0.0)
	)

	return &SingleVarLinEq{a, b, isHorizontal}, nil
}

// YAt returns the value of the single variable at a given x value.
func (svle SingleVarLinEq) YAt(x float64) float64 {
	return svle.a*x + svle.b
}

// XAt returns the value of the single variable at a given y value.
// When the slope is zero (horizontal line), the function returns an error.
func (svle SingleVarLinEq) XAt(y float64) (float64, error) {
	if svle.isHorizontal {
		return 0.0, fmt.Errorf("slope is zero")
	}

	return (y - svle.b) / svle.a, nil
}

func (svle SingleVarLinEq) Equals(other SingleVarLinEq) bool {
	return nums.FloatsEqual(svle.a, other.a) &&
		nums.FloatsEqual(svle.b, other.b)
}
