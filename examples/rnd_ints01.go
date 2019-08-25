// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/rnd"
	"github.com/dicksontsai/gosl/utl"
)

func main() {

	// initialise seed with fixed number; use 0 to use current time
	rnd.Init(1234)

	// allocate slice for integers
	nsamples := 100000
	nints := 10
	vals := make([]int, nsamples)

	// using the rnd.Int function
	t0 := time.Now()
	for i := 0; i < nsamples; i++ {
		vals[i] = rnd.Int(0, nints-1)
	}
	io.Pf("time elapsed = %v\n", time.Now().Sub(t0))

	// text histogram
	hist := rnd.IntHistogram{Stations: utl.IntRange(nints + 1)}
	hist.Count(vals, true)
	io.Pf(rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))

	// using the rnd.Ints function
	t0 = time.Now()
	rnd.Ints(vals, 0, nints-1)
	io.Pf("time elapsed = %v\n", time.Now().Sub(t0))

	// text histogram
	hist.Count(vals, true)
	io.Pf(rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))
}
