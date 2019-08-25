// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rw

import "github.com/dicksontsai/gosl/io"

func atob(s string) bool {
	if s == ".t." || s == ".T." {
		return true
	}
	if s == ".f." || s == ".F." {
		return false
	}
	return io.Atob(s)
}
