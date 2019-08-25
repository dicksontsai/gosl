// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/plt"
	"github.com/dicksontsai/gosl/utl"
)

// TwoVarsFuncType defines a function to plot contours (len(x)==2)
type TwoVarsFuncType func(x []float64) float64

// PlotTwoVarsContour plots contour for two variables problem. len(x) == 2
//  Input
//   x       -- solution. can be <nil>
//   np      -- number of points for contour
//   extra   -- called just before saving figure
//   axequal -- axis.equal
//   vmin    -- min 0 values
//   vmax    -- max 1 values
//   f       -- function to plot filled contour. can be <nil>
//   gs      -- functions to plot contour @ level 0. can be <nil>
//   argsF   -- plot arguments for f function. can be nil
//   argsG   -- plot arguments for gs functions. can be nil
func PlotTwoVarsContour(x []float64, np int, extra func(), axequal bool,
	vmin, vmax []float64, argsF, argsG *plt.A, f TwoVarsFuncType, gs ...TwoVarsFuncType) {
	chk.IntAssert(len(vmin), 2)
	chk.IntAssert(len(vmax), 2)
	V0, V1 := utl.MeshGrid2d(vmin[0], vmax[0], vmin[1], vmax[1], np, np)
	var Zf [][]float64
	var Zg [][][]float64
	if f != nil {
		Zf = utl.Alloc(np, np)
	}
	if len(gs) > 0 {
		Zg = utl.Deep3alloc(len(gs), np, np)
	}
	xtmp := make([]float64, 2)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			xtmp[0], xtmp[1] = V0[i][j], V1[i][j]
			if f != nil {
				Zf[i][j] = f(xtmp)
			}
			for k, g := range gs {
				Zg[k][i][j] = g(xtmp)
			}
		}
	}
	if f != nil {
		plt.ContourF(V0, V1, Zf, argsF)
	}
	for k := range gs {
		if argsG == nil {
			argsG = &plt.A{Levels: []float64{0}, Colors: []string{"yellow"}, Lw: 2}
		}
		plt.ContourL(V0, V1, Zg[k], argsG)
	}
	if x != nil {
		plt.PlotOne(x[0], x[1], &plt.A{C: "r", Ls: "none", M: "*", L: "optimum", Z: 10})
	}
	if extra != nil {
		extra()
	}
	plt.Cross(0, 0, &plt.A{C: "grey"})
	plt.SetXnticks(11)
	plt.SetYnticks(11)
	if axequal {
		plt.Equal()
	}
}
