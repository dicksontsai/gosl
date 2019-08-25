// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"math"

	"github.com/dicksontsai/gosl/fun"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/la"
	"github.com/dicksontsai/gosl/plt"
	"github.com/dicksontsai/gosl/utl"
)

// History holds history of optmization using directiors; e.g. for Debugging
type History struct {

	// data
	Ndim  int         // dimension of x-vector
	HistX []la.Vector // [it] history of x-values (position)
	HistU []la.Vector // [it] history of u-values (direction)
	HistF []float64   // [it] history of f-values
	HistI []float64   // [it] index of iteration

	// configuration
	NptsI   int       // number of points for contour
	NptsJ   int       // number of points for contour
	RangeXi []float64 // {ximin, ximax} [may be nil for default]
	RangeXj []float64 // {xjmin, xjmax} [may be nil for default]
	GapXi   float64   // expand {ximin, ximax}
	GapXj   float64   // expand {ximin, ximax}

	// internal
	ffcn fun.Sv // f({x}) function
}

// NewHistory returns new object
func NewHistory(nMaxIt int, f0 float64, x0 la.Vector, ffcn fun.Sv) (o *History) {
	o = new(History)
	o.Ndim = len(x0)
	o.HistX = append(o.HistX, x0.GetCopy())
	o.HistU = append(o.HistU, nil)
	o.HistF = append(o.HistF, f0)
	o.HistI = append(o.HistI, 0)
	o.NptsI = 41
	o.NptsJ = 41
	o.ffcn = ffcn
	return
}

// Append appends new x and u vectors, and updates F and I arrays
func (o *History) Append(fx float64, x, u la.Vector) {
	o.HistX = append(o.HistX, x.GetCopy())
	o.HistU = append(o.HistU, u.GetCopy())
	o.HistF = append(o.HistF, fx)
	o.HistI = append(o.HistI, float64(len(o.HistI)))
}

// Limits computes range of X variables
func (o *History) Limits() (Xmin []float64, Xmax []float64) {
	Xmin = make([]float64, o.Ndim)
	Xmax = make([]float64, o.Ndim)
	for j := 0; j < o.Ndim; j++ {
		Xmin[j] = math.MaxFloat64
		Xmax[j] = math.SmallestNonzeroFloat64
		for _, x := range o.HistX {
			Xmin[j] = utl.Min(Xmin[j], x[j])
			Xmax[j] = utl.Max(Xmax[j], x[j])
		}
	}
	return
}

// PlotC plots contour
func (o *History) PlotC(iDim, jDim int, xref la.Vector) {

	// limits
	var Xmin, Xmax []float64
	if len(o.RangeXi) != 2 || len(o.RangeXj) != 2 {
		Xmin, Xmax = o.Limits()
	}

	// i-range
	var ximin, ximax float64
	if len(o.RangeXi) == 2 {
		ximin, ximax = o.RangeXi[0], o.RangeXi[1]
	} else {
		ximin, ximax = Xmin[iDim], Xmax[iDim]
	}

	// j-range
	var xjmin, xjmax float64
	if len(o.RangeXj) == 2 {
		xjmin, xjmax = o.RangeXj[0], o.RangeXj[1]
	} else {
		xjmin, xjmax = Xmin[jDim], Xmax[jDim]
	}

	// use gap
	ximin -= o.GapXi
	ximax += o.GapXi
	xjmin -= o.GapXj
	xjmax += o.GapXj

	// contour
	xvec := xref.GetCopy()
	xx, yy, zz := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, o.NptsI, o.NptsJ, func(r, s float64) float64 {
		xvec[iDim], xvec[jDim] = r, s
		return o.ffcn(xvec)
	})
	plt.ContourF(xx, yy, zz, nil)

	// labels
	plt.SetLabels(io.Sf("$x_{%d}$", iDim), io.Sf("$x_{%d}$", jDim), nil)
}

// PlotX plots trajectory of x-points
func (o *History) PlotX(iDim, jDim int, xref la.Vector, argsArrow *plt.A) {

	// trajectory
	x2d := la.NewVector(2)
	u2d := la.NewVector(2)
	for k := 0; k < len(o.HistX)-1; k++ {
		x := o.HistX[k]
		u := o.HistU[1+k]
		x2d[0], x2d[1] = x[iDim], x[jDim]
		u2d[0], u2d[1] = u[iDim], u[jDim]
		if u.Norm() > 1e-10 {
			plt.PlotOne(x2d[0], x2d[1], &plt.A{C: "y", M: "o", Z: 10, NoClip: true})
			plt.DrawArrow2d(x2d, u2d, false, 1, argsArrow)
		}
	}

	// final point
	l := len(o.HistX) - 1
	plt.PlotOne(o.HistX[l][iDim], o.HistX[l][jDim], &plt.A{C: "y", M: "*", Ms: 10, Z: 10, NoClip: true})

	// labels
	plt.SetLabels(io.Sf("$x_{%d}$", iDim), io.Sf("$x_{%d}$", jDim), nil)
}

// PlotF plots convergence on F values versus iteration numbers
func (o *History) PlotF(args *plt.A) {
	if args == nil {
		args = &plt.A{C: plt.C(2, 0), M: ".", Ls: "-", Lw: 2, NoClip: true}
	}
	l := len(o.HistI) - 1
	plt.Plot(o.HistI, o.HistF, args)
	plt.Text(o.HistI[0], o.HistF[0], io.Sf("%.3f", o.HistF[0]), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "left", Va: "top", NoClip: true})
	plt.Text(o.HistI[l], o.HistF[l], io.Sf("%.3f", o.HistF[l]), &plt.A{C: plt.C(0, 0), Fsz: 7, Ha: "right", Va: "bottom", NoClip: true})
	plt.Gll("$iteration$", "$f(x)$", nil)
	plt.HideTRborders()
}

// PlotAll2d plots contour using PlotC, trajectory using PlotX, and convergence on F values using
// PlotF for history data with ndim >= 2
func (o *History) PlotAll2d(name string, xref la.Vector) {

	clr := "orange"

	argsArrow := &plt.A{C: clr, Scale: 40}
	argsF := &plt.A{C: clr, Lw: 3, L: name, NoClip: true}

	o.GapXi = 0.1
	o.GapXj = 0.1

	plt.SplotGap(0.25, 0.25)

	plt.Subplot(2, 1, 1)
	o.PlotC(0, 1, xref)
	o.PlotX(0, 1, xref, argsArrow)

	plt.Subplot(2, 1, 2)
	o.PlotF(argsF)
}

// PlotAll3d plots contour using PlotC, trajectory using PlotX, and convergence on F values using
// PlotF for history data with ndim >= 3
func (o *History) PlotAll3d(name string, xref la.Vector) {

	clr := "orange"

	argsArrow := &plt.A{C: clr, Scale: 40}
	argsF := &plt.A{C: clr, Lw: 3, L: name, NoClip: true}

	o.GapXi = 0.1
	o.GapXj = 0.1

	plt.SplotGap(0.25, 0.25)

	plt.Subplot(2, 2, 1)
	o.PlotC(0, 1, xref)
	o.PlotX(0, 1, xref, argsArrow)

	plt.Subplot(2, 2, 2)
	o.PlotC(1, 2, xref)
	o.PlotX(1, 2, xref, argsArrow)

	plt.Subplot(2, 2, 3)
	o.PlotC(2, 0, xref)
	o.PlotX(2, 0, xref, argsArrow)

	plt.Subplot(2, 2, 4)
	o.PlotF(argsF)
}

// CompareHistory2d generate plots to compare two history data with ndim >= 2
func CompareHistory2d(name1, name2 string, hist1, hist2 *History, xref1, xref2 la.Vector) {

	clr1 := "orange"
	clr2 := "#5a5252"

	argsArrow1 := &plt.A{C: clr1, Scale: 40}
	argsArrow2 := &plt.A{C: clr2, Scale: 10}
	argsF1 := &plt.A{C: clr1, Lw: 5, L: name1, NoClip: true}
	argsF2 := &plt.A{C: clr2, Lw: 2, L: name2, NoClip: true}

	Xmin1, Xmax1 := hist1.Limits()
	Xmin2, Xmax2 := hist2.Limits()
	hist1.RangeXi = make([]float64, 2)
	hist1.RangeXj = make([]float64, 2)

	hist1.GapXi = 0.1
	hist1.GapXj = 0.1

	plt.SplotGap(0.25, 0.25)

	plt.Subplot(2, 1, 1)
	hist1.RangeXi[0] = utl.Min(Xmin1[0], Xmin2[0])
	hist1.RangeXi[1] = utl.Max(Xmax1[0], Xmax2[0])
	hist1.RangeXj[0] = utl.Min(Xmin1[1], Xmin2[1])
	hist1.RangeXj[1] = utl.Max(Xmax1[1], Xmax2[1])
	hist1.PlotC(0, 1, xref1)
	hist1.PlotX(0, 1, xref1, argsArrow1)
	hist2.PlotX(0, 1, xref2, argsArrow2)

	plt.Subplot(2, 1, 2)
	hist1.PlotF(argsF1)
	hist2.PlotF(argsF2)
}

// CompareHistory3d generate plots to compare two history data with ndim >= 3
func CompareHistory3d(name1, name2 string, hist1, hist2 *History, xref1, xref2 la.Vector) {

	clr1 := "orange"
	clr2 := "#5a5252"

	argsArrow1 := &plt.A{C: clr1, Scale: 40}
	argsArrow2 := &plt.A{C: clr2, Scale: 10}
	argsF1 := &plt.A{C: clr1, Lw: 5, L: name1, NoClip: true}
	argsF2 := &plt.A{C: clr2, Lw: 2, L: name2, NoClip: true}

	Xmin1, Xmax1 := hist1.Limits()
	Xmin2, Xmax2 := hist2.Limits()
	hist1.RangeXi = make([]float64, 2)
	hist1.RangeXj = make([]float64, 2)

	hist1.GapXi = 0.1
	hist1.GapXj = 0.1

	plt.SplotGap(0.25, 0.25)

	plt.Subplot(2, 2, 1)
	hist1.RangeXi[0] = utl.Min(Xmin1[0], Xmin2[0])
	hist1.RangeXi[1] = utl.Max(Xmax1[0], Xmax2[0])
	hist1.RangeXj[0] = utl.Min(Xmin1[1], Xmin2[1])
	hist1.RangeXj[1] = utl.Max(Xmax1[1], Xmax2[1])
	hist1.PlotC(0, 1, xref1)
	hist1.PlotX(0, 1, xref1, argsArrow1)
	hist2.PlotX(0, 1, xref2, argsArrow2)

	plt.Subplot(2, 2, 2)
	hist1.RangeXi[0] = utl.Min(Xmin1[1], Xmin2[1])
	hist1.RangeXi[1] = utl.Max(Xmax1[1], Xmax2[1])
	hist1.RangeXj[0] = utl.Min(Xmin1[2], Xmin2[2])
	hist1.RangeXj[1] = utl.Max(Xmax1[2], Xmax2[2])
	hist1.PlotC(1, 2, xref1)
	hist1.PlotX(1, 2, xref1, argsArrow1)
	hist2.PlotX(1, 2, xref2, argsArrow2)

	plt.Subplot(2, 2, 3)
	hist1.RangeXi[0] = utl.Min(Xmin1[2], Xmin2[2])
	hist1.RangeXi[1] = utl.Max(Xmax1[2], Xmax2[2])
	hist1.RangeXj[0] = utl.Min(Xmin1[0], Xmin2[0])
	hist1.RangeXj[1] = utl.Max(Xmax1[0], Xmax2[0])
	hist1.PlotC(2, 0, xref1)
	hist1.PlotX(2, 0, xref1, argsArrow1)
	hist2.PlotX(2, 0, xref2, argsArrow2)

	plt.Subplot(2, 2, 4)
	hist1.PlotF(argsF1)
	hist2.PlotF(argsF2)
}
