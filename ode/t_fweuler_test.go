// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"testing"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/plt"
)

func TestFwEuler01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FwEuler01. Forward-Euler")

	// problem
	p := ProbHwEq11()

	// configuration
	conf := NewConfig("fweuler", "", nil)
	conf.SetFixedH(p.Dx, p.Xf)
	conf.SetStepOut(true, nil)

	// solver
	sol := NewSolver(p.Ndim, conf, p.Fcn, p.Jac, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(p.Y, 0.0, p.Xf)

	// check Stat
	chk.Int(tst, "number of F evaluations ", sol.Stat.Nfeval, 40)
	chk.Int(tst, "number of J evaluations ", sol.Stat.Njeval, 0)
	chk.Int(tst, "total number of steps   ", sol.Stat.Nsteps, 40)
	chk.Int(tst, "number of accepted steps", sol.Stat.Naccepted, 0)
	chk.Int(tst, "number of rejected steps", sol.Stat.Nrejected, 0)
	chk.Int(tst, "number of decompositions", sol.Stat.Ndecomp, 0)
	chk.Int(tst, "number of lin solutions ", sol.Stat.Nlinsol, 0)
	chk.Int(tst, "max number of iterations", sol.Stat.Nitmax, 0)

	// check results
	chk.Float64(tst, "yFin", 0.004753, p.Y[0], p.CalcYana(0, p.Xf))

	// plot
	if chk.Verbose {
		plt.Reset(true, nil)
		p.Plot("FwEuler", 0, sol.Out, 101, true, nil, nil)
		plt.Save("/tmp/gosl/ode", "fweuler01")
	}
}
