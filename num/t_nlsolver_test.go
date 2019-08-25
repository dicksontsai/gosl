// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/la"
)

func TestNls01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nls01. 2 eqs system")

	ffcn := func(fx, x la.Vector) {
		fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
		fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
		return
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
		dfdx.Start()
		dfdx.Put(0, 0, 3.0*x[0]*x[0])
		dfdx.Put(0, 1, 1.0)
		dfdx.Put(1, 0, -1.0)
		dfdx.Put(1, 1, 3.0*x[1]*x[1])
		return
	}

	x := []float64{0.5, 0.5}
	atol, rtol, ftol := 1e-10, 1e-10, 10*MACHEPS
	fx := make([]float64, len(x))
	neq := len(x)

	prms := map[string]float64{
		"atol":      atol,
		"rtol":      rtol,
		"ftol":      ftol,
		"linSearch": 1.0,
	}

	io.PfYel("\n-------------------- Analytical Jacobian -------------------\n")

	// init
	var nlsAna NlSolver
	nlsAna.Init(neq, ffcn, Jfcn, nil, false, false, prms)
	defer nlsAna.Free()

	// solve
	nlsAna.Solve(x, false)

	// check
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{1.0, 0.0})
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "f(x) = 0? ", 1e-16, fx, []float64{})

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	nlsAna.CheckJ(x, 1e-5, false, true)

	io.PfYel("\n\n-------------------- Numerical Jacobian --------------------\n")
	xx := []float64{0.5, 0.5}

	// init
	var nlsNum NlSolver
	nlsNum.Init(neq, ffcn, nil, nil, false, true, prms)
	defer nlsNum.Free()

	// solve
	nlsNum.Solve(xx, false)

	// check
	ffcn(fx, xx)
	io.Pf("xx    = %v  expected = %v\n", xx, []float64{1.0, 0.0})
	io.Pf("f(xx) = %v\n", fx)
	chk.Array(tst, "f(x) = 0? ", 1e-16, fx, []float64{})
	chk.Array(tst, "x == xx", 1e-15, x, xx)

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	nlsAna.CheckJ(x, 1e-5, false, true)
}

func TestNls02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nls02. 2 eqs system with exp function")

	ffcn := func(fx, x la.Vector) {
		fx[0] = 2.0*x[0] - x[1] - math.Exp(-x[0])
		fx[1] = -x[0] + 2.0*x[1] - math.Exp(-x[1])
		return
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
		dfdx.Start()
		dfdx.Put(0, 0, 2.0+math.Exp(-x[0]))
		dfdx.Put(0, 1, -1.0)
		dfdx.Put(1, 0, -1.0)
		dfdx.Put(1, 1, 2.0+math.Exp(-x[1]))
		return
	}

	x := []float64{5.0, 5.0}
	atol, rtol, ftol := 1e-10, 1e-10, 10*MACHEPS
	fx := make([]float64, len(x))
	neq := len(x)

	prms := map[string]float64{
		"atol":      atol,
		"rtol":      rtol,
		"ftol":      ftol,
		"linSearch": 1.0,
	}

	io.PfYel("\n-------------------- Analytical Jacobian -------------------\n")

	// init
	var nlsAna NlSolver
	nlsAna.Init(neq, ffcn, Jfcn, nil, false, false, prms)
	defer nlsAna.Free()

	// solve
	nlsAna.Solve(x, false)

	// check
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5671, 0.5671})
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "f(x) = 0? ", 1e-14, fx, []float64{})

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	nlsAna.CheckJ(x, 1e-5, false, true)

	io.PfYel("\n\n-------------------- Numerical Jacobian --------------------\n")
	xx := []float64{5.0, 5.0}

	// init
	var nlsNum NlSolver
	nlsNum.Init(neq, ffcn, nil, nil, false, true, prms)
	defer nlsNum.Free()

	// solve
	nlsNum.Solve(xx, false)

	// check
	ffcn(fx, x)
	io.Pf("xx    = %v  expected = %v\n", x, []float64{0.5671, 0.5671})
	io.Pf("f(xx) = %v\n", fx)
	chk.Array(tst, "f(x) = 0? ", 1e-14, fx, []float64{})
	chk.Array(tst, "x == xx", 1e-15, x, xx)

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	nlsAna.CheckJ(x, 1e-5, false, true)
}

func TestNls03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nls03. 2 eqs system with trig functions")

	e := math.E
	ffcn := func(fx, x la.Vector) {
		fx[0] = 0.5*sin(x[0]*x[1]) - 0.25*x[1]/pi - 0.5*x[0]
		fx[1] = (1.0-0.25/pi)*(math.Exp(2.0*x[0])-e) + e*x[1]/pi - 2.0*e*x[0]
		return
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
		dfdx.Start()
		dfdx.Put(0, 0, 0.5*x[1]*cos(x[0]*x[1])-0.5)
		dfdx.Put(0, 1, 0.5*x[0]*cos(x[0]*x[1])-0.25/pi)
		dfdx.Put(1, 0, (2.0-0.5/pi)*math.Exp(2.0*x[0])-2.0*e)
		dfdx.Put(1, 1, e/pi)
		return
	}
	JfcnD := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 0.5*x[1]*cos(x[0]*x[1])-0.5)
		dfdx.Set(0, 1, 0.5*x[0]*cos(x[0]*x[1])-0.25/pi)
		dfdx.Set(1, 0, (2.0-0.5/pi)*math.Exp(2.0*x[0])-2.0*e)
		dfdx.Set(1, 1, e/pi)
		return
	}

	x := []float64{0.4, 3.0}
	fx := make([]float64, len(x))
	atol := 1e-6
	rtol := 1e-3
	ftol := 10 * MACHEPS
	neq := len(x)

	prms := map[string]float64{
		"atol":      atol,
		"rtol":      rtol,
		"ftol":      ftol,
		"linSearch": 0.0, // does not work with line search
	}

	// init
	var nlsSps NlSolver // sparse
	var nlsDen NlSolver // dense
	nlsSps.Init(neq, ffcn, Jfcn, nil, false, false, prms)
	nlsDen.Init(neq, ffcn, nil, JfcnD, true, false, prms)
	defer nlsSps.Free()
	defer nlsDen.Free()

	io.PfMag("\n/////////////////////// sparse //////////////////////////////////////////\n")

	x = []float64{0.4, 3.0}
	io.PfYel("\n--- sparse ------------- with x = %v --------------\n", x)
	nlsSps.Solve(x, false)
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{-0.2605992900257, 0.6225308965998})
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "x", 1e-13, x, []float64{-0.2605992900257, 0.6225308965998})
	chk.Array(tst, "f(x) = 0? ", 1e-11, fx, nil)

	x = []float64{0.7, 4.0}
	io.PfYel("\n--- sparse ------------- with x = %v --------------\n", x)
	//rtol = 1e-2
	nlsSps.Solve(x, false)
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5000000377836, 3.1415927055406})
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "x  ", 1e-7, x, []float64{0.5000000377836, 3.1415927055406})
	chk.Array(tst, "f(x) = 0? ", 1e-7, fx, nil)

	x = []float64{1.0, 4.0}
	io.PfYel("\n--- sparse ------------- with x = %v ---------------\n", x)
	//linSearch, chkConv := false, true  // this combination fails due to divergence
	//linSearch, chkConv := false, false // this combination works but results are different
	//linSearch, chkConv := true, true   // this combination works but results are wrong => fails
	nlsSps.Solve(x, false)
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5, pi})
	io.Pf("f(x) = %v << converges to a different solution\n", fx)
	chk.Array(tst, "f(x) = 0? ", 1e-8, fx, nil)

	io.PfMag("\n/////////////////////// dense //////////////////////////////////////////\n")

	x = []float64{0.4, 3.0}
	io.PfYel("\n--- dense ------------- with x = %v --------------\n", x)
	nlsDen.Solve(x, false)
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{-0.2605992900257, 0.6225308965998})
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "x", 1e-13, x, []float64{-0.2605992900257, 0.6225308965998})
	chk.Array(tst, "f(x) = 0? ", 1e-11, fx, nil)

	x = []float64{0.7, 4.0}
	io.PfYel("\n--- dense ------------- with x = %v --------------\n", x)
	//rtol = 1e-2
	nlsDen.Solve(x, false)
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5000000377836, 3.1415927055406})
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "x  ", 1e-7, x, []float64{0.5000000377836, 3.1415927055406})
	chk.Array(tst, "f(x) = 0? ", 1e-7, fx, nil)

	x = []float64{1.0, 4.0}
	io.PfYel("\n--- dense ------------- with x = %v ---------------\n", x)
	nlsDen.Solve(x, false)
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5, pi})
	io.Pf("f(x) = %v << converges to a different solution\n", fx)
	chk.Array(tst, "f(x) = 0? ", 1e-8, fx, nil)
}
