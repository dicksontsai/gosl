// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/dicksontsai/gosl/gm"
	"github.com/dicksontsai/gosl/plt"
)

func main() {

	// curve
	xc, yc, r := 0.5, 0.5, 1.5
	curve := gm.FactoryNurbs.Curve2dCircle(xc, yc, r)

	// configuration
	argsIdsA := &plt.A{C: "k", Fsz: 7}
	argsCtrlA := &plt.A{C: "k", M: ".", Ls: "--", L: "control"}
	argsElemsA := &plt.A{C: "b", L: "curve"}

	// plot
	ndim := 2
	npts := 41
	plt.Reset(true, &plt.A{WidthPt: 400})
	curve.DrawCtrl(ndim, true, argsCtrlA, argsIdsA)
	curve.DrawElems(ndim, npts, true, argsElemsA, nil)
	plt.HideAllBorders()
	plt.Equal()
	plt.AxisRange(-2.5, 2.5, -2.5, 2.5)
	plt.Save("/tmp/gosl", "gm_nurbs02")
}
