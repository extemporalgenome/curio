// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curio

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestMatch(t *testing.T) {
	const (
		sc = "abc \t\ndef !@#\f(ghi) jkl"
		pc = "\x00abc\x00\x01\x00\x01\f(ghi) \x01\x00"
	)
	var (
		mx   = [3]string{"def", "!@#", "jkl"}
		mr   [3]string
		s, p io.ByteScanner = strings.NewReader(sc), strings.NewReader(pc)
	)

	t.Log("Normal")
	if n, _, err := Match(s, p, mr[:]); err != io.EOF {
		t.Fatalf("err: %v != EOF", err)
	} else if n != len(sc) {
		t.Fatalf("n: %d != %d", n, len(sc))
	} else if mr != mx {
		t.Fatalf("m: %#v != %#v", mr, mx)
	}
	mx[0], mx[2] = mx[2], mx[0]
	s = NewRevByteScanner(strings.NewReader(sc), int64(len(sc)-1))
	p = NewRevByteScanner(strings.NewReader(pc), int64(len(pc)-1))

	t.Log("Reverse")
	if n, _, err := Match(s, p, mr[:]); err != io.EOF {
		t.Fatalf("err: %v != EOF", err)
	} else if n != len(sc) {
		t.Fatalf("n: %d != %d", n, len(sc))
	} else if mr != mx {
		t.Fatalf("m: %#v != %#v", mr, mx)
	}
}

func ExampleMatch() {
	s := strings.NewReader("abc\t(def) \nghi")
	p := strings.NewReader("abc\x00(\x01)\x00ghi")
	var m [1]string
	if _, _, err := Match(s, p, m[:]); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m[0])
	}
	// Output: def
}

func ExampleMatch_reverse() {
	s := strings.NewReader("abc\tdef \nghi jkl")
	p := strings.NewReader("\x00def\x00\x01\x00jkl")
	rs := NewRevByteScanner(s, int64(s.Len()-1))
	rp := NewRevByteScanner(p, int64(p.Len()-1))

	var m [1]string
	if _, _, err := Match(rs, rp, m[:]); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m[0])
	}
	// Output: ghi
}

func ExampleMatch_numeric() {
	s := strings.NewReader("zQFf837")
	p := strings.NewReader("\xd3\xd6\xd0\xd8") // all non-greedy

	var m [4]string
	if _, _, err := Match(s, p, m[:]); err != nil && err != io.EOF {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}
	// Output: [zQ Ff 8 37]
}
