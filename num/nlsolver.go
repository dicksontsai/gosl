// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/fun"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/la"
	"github.com/dicksontsai/gosl/utl"
)

// NlSolver implements a solver to nonlinear systems of equations
//   References:
//    [1] G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
//        computations. M., Mir, 1980, p.180 of the Russian edition
type NlSolver struct {

	// constants
	cteJac      bool    // constant Jacobian (Modified Newton's method)
	linSearch   bool    // use linear search
	linSchMaxIt int     // line search maximum iterations
	maxIt       int     // Newton's method maximum iterations
	chkConv     bool    // check convergence
	atol        float64 // absolute tolerance
	rtol        float64 // relative tolerance
	ftol        float64 // minimum value of fx
	fnewt       float64 // [derived] Newton's method tolerance

	// auxiliary data
	neq   int       // number of equations
	scal  la.Vector // scaling vector
	fx    la.Vector // f(x)
	mdx   la.Vector // - delta x
	useDn bool      // use dense solver (matrix inversion) instead of Umfpack (sparse)
	numJ  bool      // use numerical Jacobian (with sparse solver)

	// callbacks
	Ffcn   fun.Vv // f(x) function f:vector, x:vector
	JfcnSp fun.Tv // J(x)=dfdx Jacobian for sparse solver
	JfcnDn fun.Mv // J(x)=dfdx Jacobian for dense solver

	// output callback
	Out func(x []float64) // output callback function

	// data for Umfpack (sparse)
	Jtri    la.Triplet // triplet
	w       la.Vector  // workspace
	lis     la.Umfpack // linear solver
	lsReady bool       // linear solver is lsReady

	// data for dense solver (matrix inversion)
	J  *la.Matrix // dense Jacobian matrix
	Ji *la.Matrix // inverse of Jacobian matrix

	// data for line-search
	φ    float64
	dφdx la.Vector
	x0   la.Vector

	// stat data
	It     int // number of iterations from the last call to Solve
	NFeval int // number of calls to Ffcn (function evaluations)
	NJeval int // number of calls to Jfcn (Jacobian evaluations)
}

// Init initialises solver
//  Input:
//   useSp -- Use sparse solver with JfcnSp
//   useDn -- Use dense solver (matrix inversion) with JfcnDn
//   numJ  -- Use numeric Jacobian (sparse version only)
//   prms  -- control parameters (default values)
//             "cteJac"      = -1 [false]  constant Jacobian (Modified Newton's method)
//             "linSearch"   = -1 [false]  use linear search
//             "linSchMaxIt" = 20          linear solver maximum iterations
//             "maxIt"       = 20          Newton's method maximum iterations
//             "chkConv"     = -1 [false]  check convergence
//             "atol"        = 1e-8        absolute tolerance
//             "rtol"        = 1e-8        relative tolerance
//             "ftol"        = 1e-9        minimum value of fx
func (o *NlSolver) Init(neq int, Ffcn fun.Vv, JfcnSp fun.Tv, JfcnDn fun.Mv, useDn, numJ bool, prms map[string]float64) {

	// set default values
	o.cteJac = false
	o.linSearch = false
	o.linSchMaxIt = 20
	o.maxIt = 20
	o.chkConv = false
	atol := 1e-8
	rtol := 1e-8
	ftol := 1e-9

	// read parameters
	for k, v := range prms {
		switch k {
		case "cteJac":
			o.cteJac = v > 0
		case "linSearch":
			o.linSearch = v > 0
		case "linSchMaxIt":
			o.linSchMaxIt = int(v)
		case "maxIt":
			o.maxIt = int(v)
		case "chkConv":
			o.chkConv = v > 0
		case "atol":
			atol = v
		case "rtol":
			rtol = v
		case "ftol":
			ftol = v
		default:
			chk.Panic("parameter named %q is invalid\n", k)
		}
	}

	// set tolerances
	o.SetTols(atol, rtol, ftol, MACHEPS)

	// auxiliary data
	o.neq = neq
	o.scal = la.NewVector(o.neq)
	o.fx = la.NewVector(o.neq)
	o.mdx = la.NewVector(o.neq)

	// callbacks
	o.Ffcn, o.JfcnSp, o.JfcnDn = Ffcn, JfcnSp, JfcnDn

	// type of linear solver and Jacobian matrix (numerical or analytical: sparse only)
	o.useDn, o.numJ = useDn, numJ

	// use dense linear solver
	if o.useDn {
		o.J = la.NewMatrix(o.neq, o.neq)
		o.Ji = la.NewMatrix(o.neq, o.neq)

		// use sparse linear solver
	} else {
		o.Jtri.Init(o.neq, o.neq, o.neq*o.neq)
		if JfcnSp == nil {
			o.numJ = true
		}
		if o.numJ {
			o.w = la.NewVector(o.neq)
		}
	}

	// allocate slices for line search
	o.dφdx = la.NewVector(o.neq)
	o.x0 = la.NewVector(o.neq)
}

// Free frees memory
func (o *NlSolver) Free() {
	if !o.useDn {
		o.lis.Free()
	}
}

// SetTols set tolerances
func (o *NlSolver) SetTols(Atol, Rtol, Ftol, ϵ float64) {
	o.atol, o.rtol, o.ftol = Atol, Rtol, Ftol
	o.fnewt = utl.Max(10.0*ϵ/Rtol, utl.Min(0.03, math.Sqrt(Rtol)))
}

// Solve solves non-linear problem f(x) == 0
func (o *NlSolver) Solve(x []float64, silent bool) {

	// compute scaling vector
	la.VecScaleAbs(o.scal, o.atol, o.rtol, x) // scal = Atol + Rtol*abs(x)

	// evaluate function @ x
	o.Ffcn(o.fx, x) // fx := f(x)
	o.NFeval, o.NJeval = 1, 0

	// show message
	if !silent {
		o.msg("", 0, 0, 0, true, false)
	}

	// iterations
	var Ldx, LdxPrev, Θ float64 // RMS norm of delta x, convergence rate
	var fxMax float64
	var nfv int
	for o.It = 0; o.It < o.maxIt; o.It++ {

		// check convergence on f(x)
		fxMax = o.fx.Largest(1.0) // den = 1.0
		if fxMax < o.ftol {
			if !silent {
				o.msg("fxMax(ini)", o.It, Ldx, fxMax, false, true)
			}
			break
		}

		// show message
		if !silent {
			o.msg("", o.It, Ldx, fxMax, false, false)
		}

		// output
		if o.Out != nil {
			o.Out(x)
		}

		// evaluate Jacobian @ x
		if o.It == 0 || !o.cteJac {
			if o.useDn {
				o.JfcnDn(o.J, x)
			} else {
				if o.numJ {
					Jacobian(&o.Jtri, o.Ffcn, x, o.fx, o.w)
					o.NFeval += o.neq
				} else {
					o.JfcnSp(&o.Jtri, x)
				}
			}
			o.NJeval++
		}

		// dense solution
		if o.useDn {

			// invert matrix
			la.MatInv(o.Ji, o.J, false)

			// solve linear system (compute mdx) and compute lin-search data
			o.φ = 0.0
			for i := 0; i < o.neq; i++ {
				o.mdx[i], o.dφdx[i] = 0.0, 0.0
				for j := 0; j < o.neq; j++ {
					o.mdx[i] += o.Ji.Get(i, j) * o.fx[j] // mdx  = inv(J) * fx
					o.dφdx[i] += o.J.Get(j, i) * o.fx[j] // dφdx = tra(J) * fx
				}
				o.φ += o.fx[i] * o.fx[i]
			}
			o.φ *= 0.5

			// sparse solution
		} else {

			// init sparse solver
			if !o.lsReady {
				symmetric, verbose := false, false
				o.lis.Init(&o.Jtri, &la.SpArgs{Symmetric: symmetric, Verbose: verbose, Ordering: "", Scaling: "", Guess: nil, Communicator: nil})
				o.lsReady = true
			}

			// factorisation (must be done for all iterations)
			o.lis.Fact()

			// solve linear system => compute mdx
			o.lis.Solve(o.mdx, o.fx, false) // mdx = inv(J) * fx   false => !sumToRoot

			// compute lin-search data
			if o.linSearch {
				o.φ = 0.5 * la.VecDot(o.fx, o.fx)
				la.SpTriMatTrVecMul(o.dφdx, &o.Jtri, o.fx) // dφdx := transpose(J) * fx
			}
		}

		// update x
		Ldx = 0.0
		for i := 0; i < o.neq; i++ {
			o.x0[i] = x[i]
			x[i] -= o.mdx[i]
			Ldx += (o.mdx[i] / o.scal[i]) * (o.mdx[i] / o.scal[i])
		}
		Ldx = math.Sqrt(Ldx / float64(o.neq))

		// calculate fx := f(x) @ update x
		o.Ffcn(o.fx, x)
		o.NFeval++

		// check convergence on f(x) => avoid line-search if converged already
		fxMax = o.fx.Largest(1.0) // den = 1.0
		if fxMax < o.ftol {
			if !silent {
				o.msg("fxMax", o.It, Ldx, fxMax, false, true)
			}
			break
		}

		// check convergence on Ldx
		if Ldx < o.fnewt {
			if !silent {
				o.msg("Ldx", o.It, Ldx, fxMax, false, true)
			}
			break
		}

		// call line-search => update x and fx
		if o.linSearch {
			nfv = LineSearch(x, o.fx, o.Ffcn, o.mdx, o.x0, o.dφdx, o.φ, o.linSchMaxIt, true)
			o.NFeval += nfv
			Ldx = 0.0
			for i := 0; i < o.neq; i++ {
				Ldx += ((x[i] - o.x0[i]) / o.scal[i]) * ((x[i] - o.x0[i]) / o.scal[i])
			}
			Ldx = math.Sqrt(Ldx / float64(o.neq))
			fxMax = o.fx.Largest(1.0) // den = 1.0
			if Ldx < o.fnewt {
				if !silent {
					o.msg("Ldx(linsrch)", o.It, Ldx, fxMax, false, true)
				}
				break
			}
		}

		// check convergence rate
		if o.It > 0 && o.chkConv {
			Θ = Ldx / LdxPrev
			if Θ > 0.99 {
				chk.Panic("solver is diverging with Θ = %g (Ldx=%g, LdxPrev=%g)", Θ, Ldx, LdxPrev)
			}
		}
		LdxPrev = Ldx
	}

	// output
	if o.Out != nil {
		o.Out(x)
	}

	// check convergence
	if o.It == o.maxIt {
		chk.Panic("cannot converge after %d iterations", o.It)
	}
	return
}

// CheckJ check Jacobian matrix
//  Ouptut: cnd -- condition number (with Frobenius norm)
func (o *NlSolver) CheckJ(x []float64, tol float64, chkJnum, silent bool) (cnd float64) {

	// Jacobian matrix
	var Jmat *la.Matrix
	if o.useDn {
		Jmat = la.NewMatrix(o.neq, o.neq)
		o.JfcnDn(Jmat, x)
	} else {
		if o.numJ {
			Jacobian(&o.Jtri, o.Ffcn, x, o.fx, o.w)
		} else {
			o.JfcnSp(&o.Jtri, x)
		}
		Jmat = o.Jtri.ToDense()
	}

	// condition number
	cnd = la.MatCondNum(Jmat, "F")
	if math.IsInf(cnd, 0) || math.IsNaN(cnd) {
		chk.Panic("condition number is Inf or NaN: %v", cnd)
	}

	// numerical Jacobian
	if !chkJnum {
		return
	}
	var Jtmp la.Triplet
	ws := la.NewVector(o.neq)
	o.Ffcn(o.fx, x)
	Jtmp.Init(o.neq, o.neq, o.neq*o.neq)
	Jacobian(&Jtmp, o.Ffcn, x, o.fx, ws)
	Jnum := Jtmp.ToMatrix(nil).ToDense()
	for i := 0; i < o.neq; i++ {
		for j := 0; j < o.neq; j++ {
			chk.PrintAnaNum(io.Sf("J[%d][%d]", i, j), tol, Jmat.Get(i, j), Jnum.Get(i, j), !silent)
		}
	}
	maxdiff := Jmat.MaxDiff(Jnum)
	if maxdiff > tol {
		chk.Panic("maxdiff = %g\n", maxdiff)
	}
	return
}

// msg prints information on residuals
func (o *NlSolver) msg(typ string, it int, Ldx, fxMax float64, first, last bool) {
	if first {
		io.Pf("\n%4s%23s%23s\n", "it", "Ldx", "fxMax")
		io.Pf("%4s%23s%23s\n", "", io.Sf("(%7.1e)", o.fnewt), io.Sf("(%7.1e)", o.ftol))
		return
	}
	io.Pf("%4d%23.15e%23.15e\n", it, Ldx, fxMax)
	if last {
		io.Pf(". . . converged with %s. nit=%d, nFeval=%d, nJeval=%d\n", typ, it, o.NFeval, o.NJeval)
	}
}
