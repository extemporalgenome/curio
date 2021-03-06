// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package curio provides io utility implementations.
package curio

import "errors"

var (
	ErrByteMismatch = errors.New("curio: Match byte mismatch")
	ErrInvalidRune  = errors.New("curio: invalid rune")
)
