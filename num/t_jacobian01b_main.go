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

func main() {

	mpi.Start()
	defer mpi.Stop()

	comm := mpi.NewCommunicator(nil)

	if comm.Rank() == 0 {
		chk.PrintTitle("TestJacobian 01b (MPI)")
	}
	if comm.Size() != 2 {
		io.Pf("this tests needs MPI 2 processors\n")
		return
	}

	ffcn := func(fx, x la.Vector) {
		fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
		fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
	}
	Jfcn := func(dfdx *la.Triplet, x la.Vector) {
		dfdx.Start()
		if false {
			if comm.Rank() == 0 {
				dfdx.Put(0, 0, 3.0*x[0]*x[0])
				dfdx.Put(1, 0, -1.0)
			} else {
				dfdx.Put(0, 1, 1.0)
				dfdx.Put(1, 1, 3.0*x[1]*x[1])
			}
		} else {
			if comm.Rank() == 0 {
				dfdx.Put(0, 0, 3.0*x[0]*x[0])
				dfdx.Put(0, 1, 1.0)
			} else {
				dfdx.Put(1, 0, -1.0)
				dfdx.Put(1, 1, 3.0*x[1]*x[1])
			}
		}
	}
	x := []float64{0.5, 0.5}
	var tst testing.T
	num.CompareJacMpi(&tst, comm, ffcn, Jfcn, x, 1e-8, true)
}
