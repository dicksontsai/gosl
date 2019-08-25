// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graph implements solvers based on Graph theory
package graph

import (
	"math"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/io"
	"github.com/dicksontsai/gosl/utl"
)

// Graph defines a graph structure
type Graph struct {

	// input
	Edges    [][]int     // [nedges][2] edges (connectivity)
	WeightsE []float64   // [nedges] weights of edges. can be <nil>
	Verts    [][]float64 // [nverts][ndim] vertices. can be <nil>
	WeightsV []float64   // [nverts] weights of vertices. can be <nil>

	// auxiliary
	Shares   map[int][]int // [nverts] edges sharing a vertex
	Key2edge map[int]int   // maps (i,j) vertex to edge index
	Dist     [][]float64   // [nverts][nverts] distances
	Next     [][]int       // [nverts][nverts] next tree connection. -1 means no connection
}

// Init initialises graph
//  Input:
//    edges    -- [nedges][2] edges (connectivity)
//    weightsE -- [nedges] weights of edges. can be <nil>
//    verts    -- [nverts][ndim] vertices. can be <nil>
//    weightsV -- [nverts] weights of vertices. can be <nil>
func (o *Graph) Init(edges [][]int, weightsE []float64, verts [][]float64, weightsV []float64) {
	o.Edges, o.WeightsE = edges, weightsE
	o.Verts, o.WeightsV = verts, weightsV
	o.Shares = make(map[int][]int)
	o.Key2edge = make(map[int]int)
	for k, edge := range o.Edges {
		i, j := edge[0], edge[1]
		utl.IntIntsMapAppend(o.Shares, i, k)
		utl.IntIntsMapAppend(o.Shares, j, k)
		o.Key2edge[o.HashEdgeKey(i, j)] = k
	}
	if o.Verts != nil {
		chk.IntAssert(len(o.Verts), len(o.Shares))
	}
	nv := len(o.Shares)
	o.Dist = utl.Alloc(nv, nv)
	o.Next = utl.IntAlloc(nv, nv)
}

// Nverts returns the number of vertices
func (o *Graph) Nverts() int {
	return len(o.Shares)
}

// GetEdge performs a lookup on Key2edge map and returs id of edge for given nodes ides
func (o *Graph) GetEdge(i, j int) (k int) {
	key := o.HashEdgeKey(i, j)
	var ok bool
	if k, ok = o.Key2edge[key]; !ok {
		chk.Panic("cannot find edge from %d to %d\n", i, j)
	}
	return
}

// ShortestPaths computes the shortest paths in a graph defined as follows
//
//          [10]
//       0 ––––––→ 3            numbers in brackets
//       |         ↑            indicate weights
//   [5] |         | [1]
//       ↓         |
//       1 ––––––→ 2
//           [3]                ∞ means that there are no
//                              connections from i to j
//   graph:  j= 0  1  2  3
//              -----------  i=
//              0  5  ∞ 10 |  0  ⇒  w(0→1)=5, w(0→3)=10
//              ∞  0  3  ∞ |  1  ⇒  w(1→2)=3
//              ∞  ∞  0  1 |  2  ⇒  w(2→3)=1
//              ∞  ∞  ∞  0 |  3
//  Input:
//   method -- FW: Floyd-Warshall method
func (o *Graph) ShortestPaths(method string) {
	if method != "FW" {
		chk.Panic("ShortestPaths works with FW (Floyd-Warshall) method only for now")
	}
	o.CalcDist()
	nv := len(o.Dist)
	var sum float64
	for k := 0; k < nv; k++ {
		for i := 0; i < nv; i++ {
			for j := 0; j < nv; j++ {
				sum = o.Dist[i][k] + o.Dist[k][j]
				if o.Dist[i][j] > sum {
					o.Dist[i][j] = sum
					o.Next[i][j] = o.Next[i][k]
				}
			}
		}
	}
	return
}

// Path returns the path from source (s) to destination (t)
//  Note: ShortestPaths method must be called first
func (o *Graph) Path(s, t int) (p []int) {
	if o.Next[s][t] < 0 {
		return
	}
	p = []int{s}
	u := s
	for u != t {
		u = o.Next[u][t]
		p = append(p, u)
	}
	return
}

// CalcDist computes distances beetween all vertices and initialises 'Next' matrix
func (o *Graph) CalcDist() {
	nv := len(o.Dist)
	for i := 0; i < nv; i++ {
		for j := 0; j < nv; j++ {
			if i == j {
				o.Dist[i][j] = 0
			} else {
				o.Dist[i][j] = math.MaxFloat64
			}
			o.Next[i][j] = -1
		}
	}
	var d float64
	for k, edge := range o.Edges {
		i, j := edge[0], edge[1]
		d = 1.0
		if o.Verts != nil {
			d = 0.0
			xa, xb := o.Verts[i], o.Verts[j]
			for dim := 0; dim < len(xa); dim++ {
				d += math.Pow(xa[dim]-xb[dim], 2.0)
			}
			d = math.Sqrt(d)
		}
		if o.WeightsE != nil {
			d *= o.WeightsE[k]
		}
		o.Dist[i][j] = d
		o.Next[i][j] = j
		if o.Dist[i][j] < 0 {
			chk.Panic("distance between vertices cannot be negative: %g\n", o.Dist[i][j])
		}
	}
	return
}

// HashEdgeKey creates a unique hash key identifying an edge
func (o *Graph) HashEdgeKey(i, j int) (edge int) {
	return i + 10000001*j
}

// StrDistMatrix returns a string representation of Dist matrix
func (o *Graph) StrDistMatrix() (l string) {
	nv := len(o.Dist)
	maxlen := 0
	for i := 0; i < nv; i++ {
		for j := 0; j < nv; j++ {
			if o.Dist[i][j] < math.MaxFloat64 {
				maxlen = utl.Imax(maxlen, len(io.Sf("%g", o.Dist[i][j])))
			}
		}
	}
	maxlen = utl.Imax(3, maxlen)
	fmts := io.Sf("%%%ds", maxlen+1)
	fmtn := io.Sf("%%%dg", maxlen+1)
	for i := 0; i < nv; i++ {
		for j := 0; j < nv; j++ {
			if o.Dist[i][j] < math.MaxFloat64 {
				l += io.Sf(fmtn, o.Dist[i][j])
			} else {
				l += io.Sf(fmts, "∞")
			}
		}
		l += "\n"
	}
	return
}

// GetAdjacency returns adjacency list as a compressed storage format for METIS
func (o *Graph) GetAdjacency() (xadj, adjncy []int32) {
	nv := o.Nverts()
	szadj := 0
	for vid := 0; vid < nv; vid++ {
		szadj += len(o.Shares[vid]) // = number of connected vertices
	}
	xadj = make([]int32, nv+1)
	adjncy = make([]int32, szadj)
	k := 0
	for vid := 0; vid < nv; vid++ {
		edges := o.Shares[vid]
		for _, eid := range edges {
			otherVid := o.Edges[eid][0]
			if otherVid == vid {
				otherVid = o.Edges[eid][1]
			}
			adjncy[k] = int32(otherVid)
			k++
		}
		xadj[1+vid] = xadj[vid] + int32(len(edges))
	}
	return
}
