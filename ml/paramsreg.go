// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"encoding/json"

	"github.com/dicksontsai/gosl/chk"
	"github.com/dicksontsai/gosl/la"
	"github.com/dicksontsai/gosl/utl"
)

// ParamsReg holds the θ and b parameters for regression computations
//
//  NOTE: Since ParamsReg is an Observable, the internal values
//        should only be changed by calling the Set... methods since
//        these methods will notify the Observers
//
type ParamsReg struct {
	utl.Observable // notifies interested parties

	// main
	theta  la.Vector // θ parameter [nFeatures]
	bias   float64   // bias parameter
	lambda float64   // regularization parameter
	degree int       // degree of polynomial

	// backup
	bkpTheta  la.Vector // copy of θ
	bkpBias   float64   // copy of b
	bkpLambda float64   // copy of λ
	bkpDegree int       // copy of degree
}

// Init initializes ParamsReg with nFeatures (number of features)
func (o *ParamsReg) Init(nFeatures int) {
	o.theta = la.NewVector(nFeatures)
	o.bkpTheta = la.NewVector(nFeatures)
}

// Backup creates an internal copy of parameters
func (o *ParamsReg) Backup() {
	copy(o.bkpTheta, o.theta)
	o.bkpBias = o.bias
	o.bkpLambda = o.lambda
	o.bkpDegree = o.degree
}

// Restore restores an internal copy of parameters and notifies observers
func (o *ParamsReg) Restore(skipNotification bool) {
	copy(o.theta, o.bkpTheta)
	o.bias = o.bkpBias
	o.lambda = o.bkpLambda
	o.degree = o.bkpDegree
	if !skipNotification {
		o.NotifyUpdate()
	}
}

// SetParams sets θ and b and notifies observers
func (o *ParamsReg) SetParams(θ la.Vector, b float64) {
	copy(o.theta, θ)
	o.bias = b
	o.NotifyUpdate()
}

// SetParam sets either θ or b (use negative indices for b). Notifies observers
//  i -- index of θ or -1 for bias
func (o *ParamsReg) SetParam(i int, value float64) {
	defer o.NotifyUpdate()
	if i < 0 {
		o.bias = value
		return
	}
	o.theta[i] = value
}

// GetParam returns either θ or b (use negative indices for b)
//  i -- index of θ or -1 for bias
func (o *ParamsReg) GetParam(i int) (value float64) {
	if i < 0 {
		return o.bias
	}
	return o.theta[i]
}

// SetThetas sets the whole vector θ and notifies observers
func (o *ParamsReg) SetThetas(θ la.Vector) {
	o.theta.Apply(1, θ)
	o.NotifyUpdate()
}

// GetThetas gets a copy of θ
func (o *ParamsReg) GetThetas() (θ la.Vector) {
	return o.theta.GetCopy()
}

// AccessThetas returns access (slice) to θ
func (o *ParamsReg) AccessThetas() (θ la.Vector) {
	return o.theta
}

// AccessBias returns access (pointer) to b
func (o *ParamsReg) AccessBias() (ptb *float64) {
	return &o.bias
}

// SetTheta sets one component of θ and notifies observers
func (o *ParamsReg) SetTheta(i int, θi float64) {
	o.theta[i] = θi
	o.NotifyUpdate()
}

// GetTheta returns the value of θ[i]
func (o *ParamsReg) GetTheta(i int) (θi float64) {
	return o.theta[i]
}

// SetBias sets b and notifies observers
func (o *ParamsReg) SetBias(b float64) {
	o.bias = b
	o.NotifyUpdate()
}

// GetBias gets a copy of b
func (o *ParamsReg) GetBias() (b float64) {
	return o.bias
}

// SetLambda sets λ and notifies observers
func (o *ParamsReg) SetLambda(λ float64) {
	o.lambda = λ
	o.NotifyUpdate()
}

// GetLambda gets a copy of λ
func (o *ParamsReg) GetLambda() (λ float64) {
	return o.lambda
}

// SetDegree sets p and notifies observers
func (o *ParamsReg) SetDegree(p int) {
	o.degree = p
	o.NotifyUpdate()
}

// GetDegree gets a copy of p
func (o *ParamsReg) GetDegree() (p int) {
	return o.degree
}

// SetJSON sets parameters from JSON string and notifies observers
func (o *ParamsReg) SetJSON(jsonString string) {

	// variable to hold input
	type jsonType struct {
		Theta  la.Vector
		Bias   float64
		Lambda float64
		Degree int
	}
	input := &jsonType{}

	// decode
	err := json.Unmarshal([]byte(jsonString), input)
	if err != nil {
		chk.Panic("cannot unmarshal json string\n")
	}

	// set data
	o.theta = input.Theta
	o.bias = input.Bias
	o.lambda = input.Lambda
	o.degree = input.Degree

	// notifications
	o.NotifyUpdate()
}
