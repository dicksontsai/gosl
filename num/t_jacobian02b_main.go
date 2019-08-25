// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"testing"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/la"
	"github.com/dicksontsai/gosl/mpi"
	"github.com/dicksontsai/gosl/num"
)

const (
	pi = math.Pi
)

func sin(x float64) float64 { return math.Sin(x) }
func cos(x float64) float64 { return math.Cos(x) }

func main() {

	mpi.Start()
	defer mpi.Stop()

	comm := mpi.NewCommunicator(nil)

	if comm.Rank() == 0 {
		chk.PrintTitle("TestJacobian 02b (MPI)")
	}
	if comm.Size() > 6 {
		io.Pf("this tests works with 6 or less MPI processors\n")
		return
	}

	ffcn := func(fx, x la.Vector) {
		fx[0] = 2.0*x[0] - x[1] + sin(x[2]) - cos(x[3]) - x[5]*x[5] - 1.0      // 0
		fx[1] = -x[0] + 2.0*x[1] + cos(x[2]) - sin(x[3]) + x[5] - 1.0          // 1
		fx[2] = x[0] + 3.0*x[1] + sin(x[3]) - cos(x[4]) - x[5]*x[5] - 1.0      // 2
		fx[3] = 2.0*x[0] + 4.0*x[1] + cos(x[3]) - cos(x[4]) + x[5] - 1.0       // 3
		fx[4] = x[0] + 5.0*x[1] - sin(x[2]) + sin(x[4]) - x[5]*x[5]*x[5] - 1.0 // 4
		fx[5] = x[0] + 6.0*x[1] - cos(x[2]) + cos(x[4]) + x[5] - 1.0           // 5
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
		dfdx.Start()
		J := [][]float64{
			{2.0, -1.0, cos(x[2]), sin(x[3]), 0.0, -2.0 * x[5]},
			{-1.0, 2.0, -sin(x[2]), -cos(x[3]), 0.0, 1.0},
			{1.0, 3.0, 0.0, cos(x[3]), sin(x[4]), -2.0 * x[5]},
			{2.0, 4.0, 0.0, -sin(x[3]), sin(x[4]), 1.0},
			{1.0, 5.0, -cos(x[2]), 0.0, cos(x[4]), -3.0 * x[5] * x[5]},
			{1.0, 6.0, sin(x[2]), 0.0, -sin(x[4]), 1.0},
		}
		id, sz, ndim := comm.Rank(), comm.Size(), 6
		start, endp1 := (id*ndim)/sz, ((id+1)*ndim)/sz
		for col := 0; col < 6; col++ {
			for row := start; row < endp1; row++ {
				dfdx.Put(row, col, J[row][col])
			}
		}
	}
	x := []float64{5.0, 5.0, pi, pi, pi, 5.0}
	var tst testing.T
	num.CompareJacMpi(&tst, comm, ffcn, Jfcn, x, 1e-6, true)
}
