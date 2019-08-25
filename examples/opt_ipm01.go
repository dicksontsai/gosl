// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/la"
	"github.com/dicksontsai/gosl/opt"
	"github.com/dicksontsai/gosl/plt"
)

func main() {

	// linear programming problem:
	//
	//   min cᵀx   s.t.   Aᵀx = b, x ≥ 0
	//    x
	//
	// specific problem:
	//
	//       min   -4*x0 - 5*x1
	//     {x0,x1}
	//
	//    s.t.  2*x0 +   x1 ≤ 3
	//            x0 + 2*x1 ≤ 3
	//          x0,x1 ≥ 0
	//
	// standard form:
	//
	//        min      -4*x0 - 5*x1
	//   {x0,x1,x2,x3}
	//
	//   s.t.
	//
	//    2*x0 +   x1 + x2     = 3
	//      x0 + 2*x1     + x3 = 3
	//    x0,x1,x2,x3 ≥ 0
	//
	// as matrix:
	//                  / x0 \
	//   [-4  -5  0  0] | x1 | = cᵀ x
	//                  | x2 |
	//                  \ x3 /
	//
	//    _            _   / x0 \
	//   |  2  1  1  0  |  | x1 | = Aᵀ x
	//   |_ 1  2  0  1 _|  | x2 |
	//                     \ x3 /
	//

	// coefficients vector
	c := []float64{-4, -5, 0, 0}

	// constraints as a sparse matrix
	var T la.Triplet
	T.Init(2, 4, 6) // 2 by 4 matrix, with 6 non-zeros
	T.Put(0, 0, 2.0)
	T.Put(0, 1, 1.0)
	T.Put(0, 2, 1.0)
	T.Put(1, 0, 1.0)
	T.Put(1, 1, 2.0)
	T.Put(1, 3, 1.0)
	Am := T.ToMatrix(nil) // compressed-column matrix

	// right-hand side
	b := []float64{3, 3}

	// solve LP
	var ipm opt.LinIpm
	defer ipm.Free()
	ipm.Init(Am, b, c, nil)
	ipm.Solve(true)

	// print solution
	io.Pf("\n")
	io.Pf("x = %v\n", ipm.X)
	io.Pf("λ = %v\n", ipm.L)
	io.Pf("s = %v\n", ipm.S)

	// check solution
	A := Am.ToDense()
	bchk := la.NewVector(2)
	la.MatVecMul(bchk, 1, A, ipm.X)
	io.Pf("b(check) = %v\n", bchk)

	// plotting
	plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
	f := func(x []float64) float64 { return c[0]*x[0] + c[1]*x[1] }
	g := func(x []float64, i int) float64 { return A.Get(i, 0)*x[0] + A.Get(i, 1)*x[1] - b[i] }
	np := 41
	argsF := &plt.A{CmapIdx: 0}
	argsG := &plt.A{Levels: []float64{0}, Colors: []string{"yellow"}, Lw: 2, Fsz: 10}
	vmin, vmax := []float64{-2.0, -2.0}, []float64{2.0, 2.0}
	opt.PlotTwoVarsContour(ipm.X[:2], np, nil, true, vmin, vmax, argsF, argsG, f,
		func(x []float64) float64 { return g(x, 0) },
		func(x []float64) float64 { return g(x, 1) },
	)
	plt.Equal()
	plt.HideAllBorders()
	plt.Gll("$x$", "$y$", &plt.A{LegOut: true})
	plt.Save("/tmp/gosl", "opt_ipm01")
}
