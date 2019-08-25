// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"
	"testing"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/utl"
)

func Test_args01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args01. arguments")

	var a A

	// plot and basic options
	a.C = "red"
	a.A = 0.5
	a.M = "o"
	a.Ls = "--"
	a.Lw = 1.2
	a.Ms = -1
	a.L = "gosl"
	a.Me = 2
	a.Z = 123
	a.Mec = "blue"
	a.Mew = 0.3
	a.Void = true
	a.NoClip = true

	// shapes
	a.Ha = "center"
	a.Va = "center"
	a.Fc = "magenta"
	a.Ec = "yellow"

	// text and extra arguments
	a.Fsz = 7

	l := a.String(false, false)
	chk.String(tst, l, "color='red',markeredgecolor='blue',markerfacecolor='none',mew=0.3,alpha=0.5,marker='o',linestyle='--',lw=1.2,label='gosl',markevery=2,zorder=123,clip_on=0,facecolor='magenta',edgecolor='yellow',ha='center',va='center',fontsize=7")
}

func Test_args02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args02. more arguments")

	a := &A{
		Colors:  []string{"red", "tan", "lime"},
		Type:    "bar",
		Stacked: true,
		NoFill:  true,
		Nbins:   10,
		Normed:  true,
	}

	l := a.String(true, false)
	chk.String(tst, l, "color=['red','tan','lime'],histtype='bar',stacked=1,fill=0,bins=10,normed=1")
}

func Test_nlevels01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("nlevels01. contour levels")

	nlevels := 3
	Z := [][]float64{
		{0, 1}, {0, 1},
		{0, 1}, {0, 1},
	}

	l := getContourLevels(nlevels, Z)
	chk.String(tst, l, "[0,0.5,1]")
}

func Test_plot01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot01. Basic")

	if chk.Verbose {

		x := utl.LinSpace(0.0, 1.0, 11)
		y := make([]float64, len(x))
		y1 := make([]float64, len(x))
		y2 := make([]float64, len(x))
		y3 := make([]float64, len(x))
		y4 := make([]float64, len(x))
		for i := 0; i < len(x); i++ {
			y[i] = x[i] * x[i]
			y1[i] = x[i]
			y2[i] = 0.5 * x[i]
			y3[i] = math.Exp(-x[i])
			y4[i] = -0.05
		}

		palette := 2
		Reset(false, nil)
		SetFontSizes(&A{Fsz: 20, FszLbl: 20, FszXtck: 10, FszYtck: 10, FontSet: "stix"})
		Plot(x, y, &A{C: C(0, palette), Lw: 3, L: "first", A: 0.5, Void: true, M: "o", Ls: "-", NoClip: true})
		Plot(y, x, &A{C: C(1, palette), Lw: 40, L: "second", M: ".", Ls: ":"})
		Plot(x, y1, &A{C: C(2, palette), Lw: 3})
		Plot(x, y2, &A{C: C(3, palette), Lw: 3})
		Plot(x, y3, &A{C: C(4, palette), Lw: 3})
		Plot(x, y4, &A{C: C(5, palette), Lw: 3})

		npts := 30
		for i := 0; i < npts; i++ {
			θ := float64(i) * math.Pi / 2.0 / float64(npts-1)
			xx := math.Cos(θ)
			yy := math.Sin(θ)
			PlotOne(xx, yy, &A{C: "k", M: M(i, 0), NoClip: true})
			xx = 0.7 * math.Cos(θ)
			yy = 0.7 * math.Sin(θ)
			PlotOne(xx, yy, &A{C: "k", M: M(i, 1), NoClip: true})
		}

		Text(0.2, 0.8, "HERE", &A{Fsz: 20, Ha: "center", Va: "center", Rot: 90})
		SetTicksX(0.1, 0.01, "%.3f")
		SetTicksY(0.2, 0.1, "%.2f")
		SetTicksRotationX(45)
		SetTicksRotationY(-45)
		SetTicksYlist(utl.LinSpace(-0.1, 1.1, 11))
		HideBorders(&A{HideR: true, HideT: true})
		//HideAllBorders()
		Gll(`$\varepsilon$`, `$\sigma$`, &A{
			LegOut:  true,
			LegNcol: 3,
			FszLeg:  14,
			HideR:   true,
		})

		Save("/tmp/gosl/plt", "t_plot01")
	}
}

func Test_plot02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot02. More Basic")

	if chk.Verbose {

		Reset(true, &A{Eps: true, WidthPt: 380})
		ReplaceAxes(-0.2, -0.1, 1.2, 1.1, 0.04, 0.04, "the x", "the y", &A{Style: "->", NoClip: true}, &A{})
		Arrow(-0.2, -0.1, 1.2, 1.1, &A{C: "orange"})
		AxHline(0, &A{C: "red"})
		AxVline(0, &A{C: "blue"})
		Annotate(0, 0, "............00", &A{C: "g", AxCoords: true})
		Annotate(1, 0, "10............", &A{C: "g", AxCoords: true, Ha: "right"})
		Annotate(1, 1, "11............", &A{C: "g", AxCoords: true, Ha: "right", Va: "top"})
		Annotate(0, 1, "............01", &A{C: "g", AxCoords: true, Va: "top"})
		Annotate(0, 0, "ann-00", &A{C: "g", FigCoords: true})
		Annotate(1, 0, "ann-10", &A{C: "g", FigCoords: true, Ha: "right"})
		Annotate(1, 1, "ann-11", &A{C: "g", FigCoords: true, Ha: "right", Va: "top"})
		Annotate(0, 1, "ann-01", &A{C: "g", FigCoords: true, Va: "top"})
		AnnotateXlabels(0, "AnnotateXlabels", &A{Fsz: 10})
		SupTitle("suptitle goes here", &A{C: "red"})
		Title("title goes here", &A{C: "cyan"})
		Text(0.5, 0.5, "TEXT", &A{C: "orange", Va: "top"})
		Text(0, 0, "TXT-00", &A{C: "m", AxCoords: true})
		Text(1, 0, "TXT-10", &A{C: "m", AxCoords: true, Ha: "right"})
		Text(1, 1, "TXT-11", &A{C: "m", AxCoords: true, Ha: "right", Va: "top"})
		Text(0, 1, "TXT-01", &A{C: "m", AxCoords: true, Va: "top"})
		//broken: Text(0, 0, "txt-00", &A{C: "m", FigCoords: true, Ha: "left", Va: "bottom"})
		//broken: Text(1, 0, "txt-10", &A{C: "m", FigCoords: true, Ha: "right", Va: "bottom"})
		//broken: Text(1, 1, "txt-11", &A{C: "m", FigCoords: true, Ha: "right", Va: "top"})
		//broken: Text(0, 1, "txt-01", &A{C: "m", FigCoords: true, Ha: "left", Va: "top"})
		Cross(0.5, 0.5, nil)
		AxisRange(-0.2, 1.3, -0.2, 1.3)
		PlotOne(0, 0, &A{M: "*"})

		Save("/tmp/gosl/plt", "t_plot02")
	}
}

func Test_plot03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot03. Contour")

	if chk.Verbose {

		// grid size
		xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

		// mesh grid
		X, Y, F := utl.MeshGrid2dF(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (z float64) {
			m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
			z = -math.Pow(m, 2.0)
			return
		})

		// configuration
		a := &A{
			NumFmt:  "%.1f",
			Lw:      0.8,
			CbarLbl: "$f(x,y)$",
			SelectC: "yellow",
			SelectV: -2.5,
			Nlevels: 10,
		}

		Reset(true, nil)
		Equal()
		ContourF(X, Y, F, a)
		SetLabels("$x$", "$y$", nil)

		Save("/tmp/gosl/plt", "t_plot03")
	}
}

func Test_plot04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot04. Contour and Quiver")

	if chk.Verbose {

		// grid size
		xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

		// mesh grid
		X, Y, F, U, V := utl.MeshGrid2dFG(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (z, u, v float64) {
			m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
			z = -math.Pow(m, 2.0)
			u = 4.0 * math.Cos(x) * math.Sin(x) * m
			v = 4.0 * math.Cos(y) * math.Sin(y) * m
			return
		})

		// configuration
		a := &A{
			Colors:   []string{"cyan", "blue", "yellow", "green"},
			Levels:   []float64{-4, -3, -2, -1, 0},
			NumFmt:   "%.1f",
			NoLines:  true,
			NoLabels: true,
			NoInline: true,
			NoCbar:   true,
			Lw:       1.5,
			SelectC:  "white",
			SelectV:  -2.5,
		}

		b := &A{
			CmapIdx: 4,
			SelectC: "black",
			SelectV: -2.5,
		}

		Reset(true, nil)
		Equal()
		ContourF(X, Y, F, a)
		ContourL(X, Y, F, b)
		Quiver(X, Y, U, V, nil)
		Grid(&A{C: "white"})

		Save("/tmp/gosl/plt", "t_plot04")
	}
}

func Test_plot05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot05. Hist")

	if chk.Verbose {

		X := [][]float64{
			{1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 4, 5, 6}, // first series
			{-1, -1, 0, 1, 2, 3},                    // second series
			{5, 6, 7, 8},                            // third series
		}

		L := []string{
			"first",
			"second",
			"third",
		}

		a := &A{
			Colors:  []string{"red", "tan", "lime"},
			Ec:      "black",
			Lw:      0.5,
			Type:    "bar",
			Stacked: true,
		}

		Reset(true, nil)
		Hist(X, L, a)
		Gll("series", "count", nil)

		Save("/tmp/gosl/plt", "t_plot05")
	}
}

func Test_plot06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot06. Plot3dLine, Plot3dPoints, Surface and Wireframe")

	if chk.Verbose {

		x := []float64{0, 1, 1, 1}
		y := []float64{0, 0, 1, 1}
		z := []float64{0, 0, 0, 1}

		np := 5

		X, Y, Z := utl.MeshGrid2dF(0, 1, 0, 1, np, np, func(x, y float64) float64 {
			return x + y
		})

		U, V, W := utl.MeshGrid2dF(0, 1, 0, 1, np, np, func(u, v float64) float64 {
			return u*u + v*v
		})

		Reset(true, nil)

		Triad(1.3, "x", "y", "z", nil, nil)

		Plot3dLine(x, y, z, nil)
		Plot3dPoints(x, y, z, nil)
		Surface(U, V, W, &A{CmapIdx: 4, Rstride: 1, Cstride: 1})

		Wireframe(X, Y, Z, &A{C: "orange", Lw: 0.4})

		elev, azim := 30.0, 20.0
		Camera(elev, azim, nil)
		AxDist(10.5)
		Scale3d(0, 1.5, 0, 1.5, 0, 1.5, true)

		Save("/tmp/gosl/plt", "t_plot06")
	}
}

func Test_plot07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot07. Triad, PlaneZ and Default3dView")

	if chk.Verbose {

		// point on plane
		p := []float64{0.5, 0.5, 0.5}

		// normal vector
		n := []float64{1, 0, 1}

		// limits and divisions
		xmin, xmax := 0.0, 1.0
		ymin, ymax := 0.0, 1.0
		nu, nv := 5, 5

		// draw
		Reset(true, nil)
		Triad(1.0, "x", "y", "z", &A{C: "orange"}, &A{C: "red"})
		PlaneZ(p, n, xmin, xmax, ymin, ymax, nu, nv, true, nil)
		Default3dView(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)

		// save
		Save("/tmp/gosl/plt", "t_plot07")
	}
}

func Test_plot08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot08. Hemisphere")

	if chk.Verbose {

		// draw
		Reset(true, nil)
		Triad(1.0, "x", "y", "z", &A{C: "orange"}, &A{C: "red"})

		alpha, height, radius := 15.0, 1.0, 0.5
		nu, nv := 7, 11
		if true {
			c1 := []float64{0.5, -0.5, 0}
			c2 := []float64{0, 1, 0}
			c3 := []float64{0, 0, 0}
			CylinderZ(c1, radius, height, nu, nv, &A{C: "grey"})
			ConeZ(c2, alpha, height, nu, nv, &A{C: "green"})
			ConeDiag(c3, alpha, height, nu, nv, nil)
			Diag3d(1, nil)
		}

		centre := []float64{0.7, 0.7, 0.0}
		radius, amin, amax := 0.3, 0.0, 180.0
		nu, nv = 21, 5
		cup := false
		Hemisphere(centre, radius, amin, amax, nu, nv, cup, &A{C: "k", Surf: true, Wire: false})

		Default3dView(-0.1, 1.1, -0.1, 1.1, -0.1, 1.1, true)

		SetLabels3d(`$x_{axis}$`, `$y_{axis}$`, `$z_{axis}$`, &A{C: "r", Fsz: 14})

		// save
		Save("/tmp/gosl/plt", "t_plot08")
	}
}

func Test_plot09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot09. Superquadric")

	if chk.Verbose {

		// draw
		Reset(true, nil)
		Triad(1.0, "$x_{axis}$", "$y_{axis}$", "$z_{axis}$", &A{C: "orange"}, &A{C: "orange"})

		// centre and radii
		c := []float64{0, 0, 0}
		r := []float64{1, 1, 1}
		a := []float64{2, 2, 2}

		// cup
		//alpMin, alpMax := -180.0, 180.0
		//etaMin, etaMax := -90.0, 0.0

		// hemisphere
		//alpMin, alpMax := -180.0, 180.0
		//etaMin, etaMax := 0.0, 90.0

		// sphere
		alpMin, alpMax := -180.0, 180.0
		etaMin, etaMax := -90.0, 90.0

		// rounded cube
		//a = []float64{10, 10, 10}

		// star
		//a = []float64{0.5, 0.5, 0.5}

		// blob
		a = []float64{2.0, 1.5, 1.0}

		// divisions
		nalp, neta := 30, 30

		// generate
		Superquadric(c, r, a, alpMin, alpMax, etaMin, etaMax, nalp, neta, &A{Surf: true})

		//Default3dView(-2.1, 2.1, -2.1, 2.1, -2.1, 2.1, true)
		Default3dView(-1.1, 1.1, -1.1, 1.1, -1.1, 1.1, true)

		// save
		Save("/tmp/gosl/plt", "t_plot09")
	}
}

func Test_plot10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot10. ZoomWindow")

	if chk.Verbose {

		// data
		x := utl.LinSpace(0.0, 100.0, 11)
		y1 := make([]float64, len(x))
		y2 := make([]float64, len(x))
		y3 := make([]float64, len(x))
		y4 := make([]float64, len(x))
		for i := 0; i < len(x); i++ {
			y1[i] = x[i] * x[i]
			y2[i] = x[i]
			y3[i] = x[i] * 100
			y4[i] = x[i] * 2
		}

		// clear figure
		Reset(false, nil)

		// plot curve on main figure
		Plot(x, y1, &A{L: "curve on old"})

		// plot curve on zoom window
		old, new := ZoomWindow(0.25, 0.5, 0.3, 0.3, nil)
		Plot(x, y2, &A{C: "r", L: "curve on new"})

		// activate main figure
		Sca(old)
		Plot(x, y3, &A{C: "orange", L: "curve ond old again"})
		Gll("x", "y", &A{LegLoc: "lower right"})

		// activate zoom window
		Sca(new)
		Plot(x, y4, &A{C: "cyan", L: "curve ond new again"})
		Gll("xnew", "ynew", nil)

		Save("/tmp/gosl/plt", "t_plot10")
	}
}

func Test_plot11(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot11. LegendX")

	if chk.Verbose {

		x := utl.LinSpace(0.0, 1.0, 11)
		y := make([]float64, len(x))
		for i := 0; i < len(x); i++ {
			y[i] = x[i] * x[i]
		}

		Reset(false, nil)
		Plot(x, y, nil)

		LegendX([]*A{
			{C: "red", M: "o", Ls: "-", Lw: 1, Ms: -1, L: "first", Me: -1},
			{C: "green", M: "s", Ls: "-", Lw: 2, Ms: 0, L: "second", Me: -1},
			{C: "blue", M: "+", Ls: "-", Lw: 3, Ms: 10, L: "third", Me: -1},
		},
			&A{LegOut: true, LegNcol: 3},
		)

		Save("/tmp/gosl/plt", "t_plot11")
	}
}

func Test_plot12(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot12. Grid2d")

	if chk.Verbose {

		X := [][]float64{
			{0, 0.5, 1}, // j=0
			{0, 0.5, 1}, // j=1
		}

		Y := [][]float64{
			{0, 0, 0}, // j=0
			{1, 1, 1}, // j=1
		}

		dx, dy := 1.1, 1.1
		U := [][]float64{
			{dx + 0.0, dx + 0.0}, // j=0
			{dx + 0.5, dx + 0.5}, // j=1
			{dx + 1.0, dx + 1.0}, // j=2
		}

		V := [][]float64{
			{dy, dy + 1}, // j=0
			{dy, dy + 1}, // j=1
			{dy, dy + 1}, // j=2
		}

		Reset(false, nil)
		Grid2d(X, Y, true, &A{C: "b", NoClip: true}, &A{C: "r", Fsz: 7})
		Grid2d(U, V, true, &A{C: "r", NoClip: true}, &A{C: "b", Fsz: 7})
		HideAllBorders()
		Equal()

		Save("/tmp/gosl/plt", "t_plot12")
	}
}

func Test_plot13(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot13. Grid3dZlevels")

	if chk.Verbose {

		X := [][]float64{
			{0, 0.5, 1}, // j=0
			{0, 0.5, 1}, // j=1
		}

		Y := [][]float64{
			{0, 0, 0}, // j=0
			{1, 1, 1}, // j=1
		}

		Zlevels := []float64{0, 1} // k=0,1

		dx, dy := 1.1, 1.1
		U := [][]float64{
			{dx + 0.0, dx + 0.0}, // j=0
			{dx + 0.5, dx + 0.5}, // j=1
			{dx + 1.0, dx + 1.0}, // j=2
		}

		V := [][]float64{
			{dy + 0, dy + 1}, // j=0
			{dy + 0, dy + 1}, // j=1
			{dy + 0, dy + 1}, // j=2
		}

		Reset(false, nil)
		Grid3dZlevels(X, Y, Zlevels, true, &A{Lw: 1, C: C(0, 0)}, &A{Fsz: 7, C: C(1, 0)})
		Grid3dZlevels(U, V, Zlevels, true, &A{Lw: 1, C: C(2, 0)}, &A{Fsz: 7, C: C(3, 0)})
		DefaultTriad(1.1)
		Default3dView(0, 1.1+dx, 0, 1.1+dy, 0, 1.1, true)

		Save("/tmp/gosl/plt", "t_plot13")
	}
}

func Test_plot14(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot14. Grid3dZlevels")

	if chk.Verbose {

		X := [][][]float64{
			{ // k=0
				{0, 0.5, 1}, // j=0
				{0, 0.5, 1}, // j=1
			},
			{ // k=1
				{0, 0.5, 1},
				{0, 0.5, 1},
			},
		}

		Y := [][][]float64{
			{
				{0, 0, 0},
				{1, 1, 1},
			},
			{
				{0, 0, 0},
				{0.5, 0.5, 0.5},
			},
		}

		Z := [][][]float64{
			{
				{0, 0, 0},
				{0, 0, 0},
			},
			{
				{1, 1, 1},
				{1, 1, 1},
			},
		}

		dx, dy := 1.1, 1.1
		U := [][][]float64{
			{ // k=0
				{dx + 0.0, dx + 0.0}, // j=0
				{dx + 0.5, dx + 0.5}, // j=1
				{dx + 1.0, dx + 1.0}, // j=2
			}, // k=1
			{
				{dx + 0.0 - 0.5, dx + 0.0 - 0.5},
				{dx + 0.5 - 0.5, dx + 0.5 - 0.5},
				{dx + 1.0 - 0.5, dx + 1.0 - 0.5},
			},
		}

		V := [][][]float64{
			{
				{dy + 0, dy + 1},
				{dy + 0, dy + 1},
				{dy + 0, dy + 1},
			},
			{
				{dy + 0, dy + 1},
				{dy + 0, dy + 1},
				{dy + 0, dy + 1},
			},
		}

		W := [][][]float64{
			{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			{
				{1, 1},
				{1, 1},
				{1, 1},
			},
		}

		dx, dy = 2, 0
		M := [][][]float64{
			{ // k=0
				{dx + 0, dx + 0}, // j=0
				{dx + 0, dx + 0}, // j=1
			},
			{
				{dx + 0.5, dx + 0.5},
				{dx + 0.5, dx + 0.5},
			},
			{
				{dx + 1, dx + 1},
				{dx + 1, dx + 1},
			},
		}

		N := [][][]float64{
			{
				{0, 1},
				{0, 1},
			},
			{
				{0, 1},
				{0, 1},
			},
			{
				{0, 1},
				{0, 1},
			},
		}

		P := [][][]float64{
			{
				{0, 0},
				{1, 1},
			},
			{
				{0, 0},
				{1, 1},
			},
			{
				{0, 0},
				{1, 1},
			},
		}

		Reset(false, nil)
		Grid3d(X, Y, Z, true, &A{Lw: 1, C: C(0, 0)}, &A{Fsz: 7, C: C(1, 0)})
		Grid3d(U, V, W, true, &A{Lw: 1, C: C(2, 0)}, &A{Fsz: 7, C: C(3, 0)})
		Grid3d(M, N, P, true, &A{Lw: 2, C: C(4, 0)}, &A{Fsz: 7, C: "k"})
		DefaultTriad(1.1)
		Default3dView(0, 1.1+dx, 0, 1.1+dy, 0, 1.1, true)

		Save("/tmp/gosl/plt", "t_plot14")
	}
}

func Test_plot15(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot15. Grid3d")

	if chk.Verbose {

		X := [][][]float64{
			{ // k=0
				{1, 2, 4, 8}, // j=0
				{1, 2, 4, 8}, // j=1
				{1, 2, 4, 8}, // j=2
			},
			{
				{2, 3, 5, 9},
				{2, 3, 5, 9},
				{2, 3, 5, 9},
			},
		}

		Y := [][][]float64{
			{
				{0, 0, 0, 0},
				{3, 3, 3, 3},
				{4, 4, 4, 4},
			},
			{
				{0, 0, 0, 0},
				{3, 3, 3, 3},
				{4, 4, 4, 4},
			},
		}

		Z := [][][]float64{
			{
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
			{
				{-2.5, -2.5, -2.5, -2.5},
				{-2.5, -2.5, -2.5, -2.5},
				{-2.5, -2.5, -2.5, -2.5},
			},
		}

		Reset(true, &A{WidthPt: 500})
		Grid3d(X, Y, Z, true, nil, nil)
		DefaultTriad(1)
		Default3dView(1, 8, 0, 4, -1, -0.5, true)

		Save("/tmp/gosl/plt", "t_plot15")
	}
}

func Test_plot16(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot14. Sphere")

	if chk.Verbose {
		Reset(true, nil)
		Sphere(nil, 1, 30, 30, &A{Surf: true})
		Default3dView(-1.1, 1.1, -1.1, 1.1, -1.1, 1.1, true)
		Save("/tmp/gosl/plt", "t_plot16")
	}
}
