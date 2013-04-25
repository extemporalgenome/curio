// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curio

import (
	"errors"
	"io"
	"unicode/utf8"
)

// NewRevByteScanner returns a backward ByteScanner.
func NewRevByteScanner(r io.ReaderAt, offset int64) io.ByteScanner {
	if offset < 0 {
		panic("negative offset is not allowed")
	}
	return &rb{r: r, o: offset + 1, p: -1}
}

const bs = 4 << 10

type rb struct {
	r    io.ReaderAt
	o    int64
	p, q int16
	u    bool
	b    [bs]byte
}

func (r *rb) ReadByte() (c byte, err error) {
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

func (r *rb) UnreadByte() error {
	if r.u {
		r.p++
		r.u = false
		return nil
	}
	return errors.New("curio: invalid use of UnreadByte")
}

var ErrInvalidRune = errors.New("curio: invalid rune")

func NewRevRuneScanner(r io.ReaderAt, offset int64) io.RuneScanner {
	if offset < 0 {
		panic("negative offset is not allowed")
	}
	return &rr{rb: rb{r: r, o: offset + 1, p: -1}, r: -1}
}

type rr struct {
	rb rb
	b  [utf8.UTFMax]byte
	r  rune
	u  bool
}

func (r *rr) ReadRune() (rn rune, size int, err error) {
	if r.u {
		r.u = false
		return r.r, utf8.RuneLen(rn), nil
	}
	for i := len(r.b) - 1; i >= 0; i-- {
		c, err := r.rb.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		r.b[i] = c
		rn, size = utf8.DecodeRune(r.b[i:])
		if rn == utf8.RuneError {
			continue
		} else if size < len(r.b)-i {
			r.r = rn
			return 0, len(r.b) - i, ErrInvalidRune
		}
		r.r = rn
		r.u = false
		return rn, size, nil
	}
	return 0, 0, ErrInvalidRune
}

func (r *rr) UnreadRune() error {
	if r.u || r.r < 0 {
		return errors.New("curio: invalid use of UnreadRune")
	}
	r.u = true
	return nil
}
