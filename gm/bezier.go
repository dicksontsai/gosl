// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/num"
	"github.com/dicksontsai/gosl/plt"
)

// BezierQuad implements a quadratic Bezier curve
//  C(t) = (1-t)² Q0  +  2 t (1-t) Q1  +  t² Q2
//       = t² A  +  2 t B  +  Q0
//  A = Q2 - 2 Q1 + Q0
//  B = Q1 - Q0
type BezierQuad struct {

	// input
	Q [][]float64 // control points; can be set outside

	// auxiliary
	P []float64 // a point on curve
}

// Point returns the x-y-z coordinates of a point on Bezier curve
func (o *BezierQuad) Point(C []float64, t float64) {
	if len(o.Q) != 3 {
		chk.Panic("Point: quadratic Bezier must be initialised first (with 3 control points)")
	}
	ndim := len(o.Q[0])
	chk.IntAssert(len(C), ndim)
	for i := 0; i < ndim; i++ {
		C[i] = (1.0-t)*(1.0-t)*o.Q[0][i] + 2.0*t*(1.0-t)*o.Q[1][i] + t*t*o.Q[2][i]
	}
	return
}

// GetPoints returns points along the curve for given parameter values
func (o *BezierQuad) GetPoints(T []float64) (X, Y, Z []float64) {
	if len(o.Q) != 3 {
		chk.Panic("GetPoints: quadratic Bezier must be initialised first (with 3 control points)")
	}
	ndim := len(o.Q[0])
	C := make([]float64, ndim)
	X = make([]float64, len(T))
	Y = make([]float64, len(T))
	if ndim > 2 {
		Z = make([]float64, len(T))
	}
	for i := 0; i < len(T); i++ {
		o.Point(C, T[i])
		X[i] = C[0]
		Y[i] = C[1]
		if ndim > 2 {
			Z[i] = C[2]
		}
	}
	return
}

// GetControlCoords returns the coordinates of control points as 1D arrays (e.g. for plotting)
func (o *BezierQuad) GetControlCoords() (X, Y, Z []float64) {
	if len(o.Q) != 3 {
		chk.Panic("GetControlCoords: quadratic Bezier must be initialised first (with 3 control points)")
	}
	ndim := len(o.Q[0])
	X = make([]float64, len(o.Q))
	Y = make([]float64, len(o.Q))
	if ndim > 2 {
		Z = make([]float64, len(o.Q))
	}
	for i := 0; i < len(o.Q); i++ {
		X[i] = o.Q[i][0]
		Y[i] = o.Q[i][1]
		if ndim > 2 {
			Z[i] = o.Q[i][2]
		}
	}
	return
}

// DistPoint returns the distance from a point to this Bezier curve
// It finds the closest projection which is stored in P
func (o *BezierQuad) DistPoint(X []float64, doplot bool) float64 {

	// TODO:
	//   1) split this into closest projections finding
	//   2) finish distance computation

	// check
	if len(o.Q) != 3 {
		chk.Panic("DistPoint: quadratic Bezier must be initialised first (with 3 control points)")
	}
	ndim := len(o.Q[0])
	chk.IntAssert(len(X), ndim)

	// solve cubic equation
	var Ai, Bi, Mi, a, b, c, d float64
	for i := 0; i < ndim; i++ {
		Ai = o.Q[2][i] - 2.0*o.Q[1][i] + o.Q[0][i]
		Bi = o.Q[1][i] - o.Q[0][i]
		Mi = o.Q[0][i] - X[i]
		a += Ai * Ai
		b += 3.0 * Ai * Bi
		c += 2.0*Bi*Bi + Mi*Ai
		d += Mi * Bi
	}
	//io.Pforan("a=%v b=%v c=%v d=%v\n", a, b, c, d)
	if math.Abs(a) < 1e-7 {
		chk.Panic("DistPoint does not yet work with this type of Bezier (straight line?):\nQ=%v\n", o.Q)
	}
	x1, x2, x3, nx := num.EqCubicSolveReal(b/a, c/a, d/a)
	io.Pfyel("\nx1=%v x2=%v x3=%v nx=%v\n", x1, x2, x3, nx)

	// auxiliary
	if len(o.P) != ndim {
		o.P = make([]float64, ndim)
	}

	// closest projections
	t := x1
	if nx == 2 {
		chk.Panic("nx=2 => not implemented yet")
	}
	if nx == 3 {
		T := []float64{x1, x2, x3}
		D := []float64{-1, -1, -1}
		ok := []bool{
			!(x1 < 0.0 || x1 > 1.0),
			!(x2 < 0.0 || x2 > 1.0),
			!(x3 < 0.0 || x3 > 1.0),
		}
		io.Pforan("ok = %v\n", ok)
		for i, t := range T {
			if ok[i] {
				o.Point(o.P, t)
				if doplot {
					plt.PlotOne(X[0], X[1], &plt.A{C: "k", M: "o"})
					plt.PlotOne(o.P[0], o.P[1], &plt.A{C: "k", M: "."})
					plt.Arrow(X[0], X[1], o.P[0], o.P[1], &plt.A{Ec: "none"})
				}
				D[i] = ppdist(X, o.P)
			}
		}
		io.Pforan("D = %v\n", D)
	}
	o.Point(o.P, t)
	io.Pfcyan("P = %v\n", o.P)
	return 0
}

// ppdist computes point-point distance
func ppdist(a, b []float64) (d float64) {
	for i := 0; i < len(a); i++ {
		d += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(d)
}
