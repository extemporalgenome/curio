// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curio

import (
	"errors"
	"io"
)

// NewRevByteScanner returns a backward ByteScanner.
func NewRevByteScanner(r io.ReaderAt, offset int64) io.ByteScanner {
	if offset < 0 {
		panic("negative offset is not allowed")
	}
	return &rev{r: r, o: offset + 1, p: -1}
}

const bs = 8 << 10

type rev struct {
	r    io.ReaderAt
	o    int64
	p, q int16
	u    bool
	b    [bs]byte
}

func (r *rev) ReadByte() (c byte, err error) {
	if r.p < 0 {
		if r.o == 0 {
			return 0, io.EOF
		}
		d := int64(bs)
		if r.o < d {
			d = r.o
		}
		r.o -= d
		r.q = int16(d)
		r.p = r.q - 1
		_, err = r.r.ReadAt(r.b[:r.q], r.o)
		if err != nil {
			return
		}
	}
	c = r.b[r.p]
	r.p--
	r.u = true
	return
}

func (r *rev) UnreadByte() error {
	if r.u {
		r.p++
		r.u = false
		return nil
	}
	return errors.New("UnreadByte: previous operation was not a read")
}
