// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"
	"time"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/plt"
)

// Hairer-Wanner VII-p2 Eq.(1.1)
func TestOde01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode01: Hairer-Wanner VII-p2 Eq.(1.1)")

	// problem
	p := ProbHwEq11()

	// FwEuler
	io.Pforan("\n. . . FwEuler . . . \n")
	_, stat1, out1 := p.Solve("fweuler", true, false)
	chk.Int(tst, "number of F evaluations ", stat1.Nfeval, 40)
	chk.Int(tst, "number of J evaluations ", stat1.Njeval, 0)
	chk.Int(tst, "total number of steps   ", stat1.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", stat1.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", stat1.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat1.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", stat1.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", stat1.Nitmax, 0)

	// BwEuler
	io.Pforan("\n. . . BwEuler . . . \n")
	_, stat2, out2 := p.Solve("bweuler", true, false)
	chk.Int(tst, "number of F evaluations ", stat2.Nfeval, 80)
	chk.Int(tst, "number of J evaluations ", stat2.Njeval, 40)
	chk.Int(tst, "total number of steps   ", stat2.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", stat2.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", stat2.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat2.Ndecomp, 40)
	chk.Int(tst, "number of lin solutions ", stat2.Nlinsol, 40)
	chk.Int(tst, "max number of iterations", stat2.Nitmax, 2)

	// MoEuler
	io.Pforan("\n. . . MoEuler . . . \n")
	_, stat3, out3 := p.Solve("moeuler", false, false)
	chk.Int(tst, "number of F evaluations ", stat3.Nfeval, 425)
	chk.Int(tst, "number of J evaluations ", stat3.Njeval, 0)
	chk.Int(tst, "total number of steps   ", stat3.Nsteps, 212)
	chk.Int(tst, "number of accepted steps", stat3.Naccepted, 212)
	chk.Int(tst, "number of rejected steps", stat3.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat3.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", stat3.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", stat3.Nitmax, 0)

	// DoPri5
	io.Pforan("\n. . . DoPri5 . . . \n")
	_, stat4, out4 := p.Solve("dopri5", false, false)
	chk.Int(tst, "number of F evaluations ", stat4.Nfeval, 242)
	chk.Int(tst, "number of J evaluations ", stat4.Njeval, 0)
	chk.Int(tst, "total number of steps   ", stat4.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", stat4.Naccepted, 40)
	chk.Int(tst, "number of rejected steps", stat4.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat4.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", stat4.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", stat4.Nitmax, 0)

	// Radau5
	io.Pforan("\n. . . Radau5 . . . \n")
	_, stat5, out5 := p.Solve("radau5", false, false)
	chk.Int(tst, "number of F evaluations ", stat5.Nfeval, 66)
	chk.Int(tst, "number of J evaluations ", stat5.Njeval, 1)
	chk.Int(tst, "total number of steps   ", stat5.Nsteps, 15)
	chk.Int(tst, "number of accepted steps", stat5.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", stat5.Nrejected, 0)
	chk.Int(tst, "number of decompositions", stat5.Ndecomp, 13)
	chk.Int(tst, "number of lin solutions ", stat5.Nlinsol, 17)
	chk.Int(tst, "max number of iterations", stat5.Nitmax, 2)

	// plot
	if chk.Verbose {
		npts := 201
		plt.Reset(true, nil)
		p.Plot("FwEuler", 0, out1, npts, true, nil, &plt.A{C: "k", M: ".", Ls: ":"})
		p.Plot("BwEuler", 0, out2, npts, false, nil, &plt.A{C: "r", M: ".", Ls: ":"})
		p.Plot("MoEuler", 0, out3, npts, false, nil, &plt.A{C: "c", M: "+", Ls: ":"})
		p.Plot("Dopri5", 0, out4, npts, false, nil, &plt.A{C: "m", M: "^", Ls: "--", Ms: 3})
		p.Plot("Radau5", 0, out5, npts, false, nil, &plt.A{C: "b", M: "s", Ls: "-", Ms: 3})
		plt.Gll("$x$", "$y$", nil)
		plt.Save("/tmp/gosl/ode", "ode1")
	}
}

// Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation
func TestOde02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode02: Hairer-Wanner VII-p5 Eq.(1.5) Van der Pol's Equation")

	// problem
	p := ProbVanDerPol(0, false)

	// configuration
	conf := NewConfig("radau5", "", nil)
	conf.SetStepOut(true, nil)

	// allocate ODE object
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve problem
	sol.Solve(p.Y, 0, p.Xf)

	// check
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2233)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 160)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 280)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 241)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 7)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 251)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 663)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T := io.ReadTable("data/vdpol_radau5_for.dat")
		X := sol.Out.GetStepX()
		for j := 0; j < p.Ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.GetStepY(j)
			plt.Subplot(p.Ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(p.Ndim+1, 1, p.Ndim+1)
		plt.Plot(X, sol.Out.GetStepH(), &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "ode2")
	}
}

// Hairer-Wanner VII-p3 Eq.(1.4) Robertson Equation
func TestOde03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode03: Hairer-Wanner VII-p3 Eq.(1.4) Robertson's Equation")

	// problem
	p := ProbRobertson()

	// configuration
	conf := NewConfig("radau5", "", nil)
	conf.SetStepOut(true, nil)

	// tolerances and initial step size
	rtol := 1e-2
	atol := rtol * 1e-6
	conf.SetTols(atol, rtol)
	conf.IniH = 1.0e-6

	// allocate ODE object
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve problem
	sol.Solve(p.Y, 0.0, p.Xf)

	// check
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 87)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 8)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 17)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 15)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 1)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 15)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 24)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 2)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.5, FszXtck: 6, FszYtck: 6})
		_, T := io.ReadTable("data/rober_radau5_cpp.dat")
		X := sol.Out.GetStepX()
		for j := 0; j < p.Ndim; j++ {
			labelA, labelB := "", ""
			if j == 2 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.GetStepY(j)
			plt.Subplot(p.Ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA, NoClip: true})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 2, Ls: "none", L: labelB, NoClip: true})
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
			plt.HideTRborders()
		}
		plt.Subplot(p.Ndim+1, 1, p.Ndim+1)
		plt.Plot(X, sol.Out.GetStepH(), &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "ode3")
	}
}

func TestOde04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ode04: Hairer-Wanner VII-p376 Transistor Amplifier")

	// problem
	p := ProbHwAmplifier()

	// configurations
	conf := NewConfig("radau5", "", nil)
	conf.SetStepOut(true, nil)
	conf.IniH = 1.0e-6 // initial step size

	// set tolerances
	atol, rtol := 1e-11, 1e-5
	conf.SetTols(atol, rtol)

	// ODE solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, p.M)
	defer sol.Free()

	// run
	t0 := time.Now()
	sol.Solve(p.Y, 0.0, p.Xf)
	io.Pfmag("elapsed time = %v\n", time.Now().Sub(t0))

	// check
	if false { // these values vary slightly in different machines
		chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 2599)
		chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 216)
		chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 275)
		chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 219)
		chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 20)
		chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 274)
		chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 792)
		chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 6)
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 450, Dpi: 150, Prop: 1.8, FszXtck: 6, FszYtck: 6})
		_, T := io.ReadTable("data/radau5_hwamplifier.dat")
		X := sol.Out.GetStepX()
		for j := 0; j < p.Ndim; j++ {
			labelA, labelB := "", ""
			if j == 4 {
				labelA, labelB = "reference", "gosl"
			}
			Yj := sol.Out.GetStepY(j)
			plt.Subplot(p.Ndim+1, 1, j+1)
			plt.Plot(T["x"], T[io.Sf("y%d", j)], &plt.A{C: "k", M: "+", L: labelA, NoClip: true})
			plt.Plot(X, Yj, &plt.A{C: "r", M: ".", Ms: 1, Ls: "none", L: labelB, NoClip: true})
			plt.AxisXmax(0.05)
			plt.HideTRborders()
			plt.Gll("$x$", io.Sf("$y_%d$", j), nil)
		}
		plt.Subplot(p.Ndim+1, 1, p.Ndim+1)
		plt.Plot(X, sol.Out.GetStepH(), &plt.A{C: "b", NoClip: true})
		plt.SetYlog()
		plt.AxisXmax(0.05)
		plt.Gll("$x$", "$\\log{(h)}$", nil)
		plt.Save("/tmp/gosl/ode", "ode4")
	}
}
