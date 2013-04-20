// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curio

import (
	"io"
	"strings"
	"testing"
)

func TestRev(t *testing.T) {
	s := strings.NewReader("abcdefghijklmnopq")
	r := NewRevByteScanner(s, int64(s.Len()-1))
	const min, max byte = 'a', 'q'
	if err := r.UnreadByte(); err == nil {
		t.Logf("expected UnreadByte error")
	}
	for b := max; ; b-- {
		c, err := r.ReadByte()
		if err == io.EOF {
			if b >= min {
				t.Fatal("early EOF")
			}
			break
		} else if err != nil {
			t.Fatal(err)
		} else if c != b {
			t.Fatalf("byte mismatch: %q != %q", c, b)
		}
		t.Logf("%q", c)
	}
	if c, err := r.ReadByte(); err != io.EOF {
		t.Fatalf("expected EOF; received %q %s", c, err)
	} else if err := r.UnreadByte(); err != nil {
		t.Fatal(err)
	} else if c, err := r.ReadByte(); err != nil {
		t.Fatal(err)
	} else if c != min {
		t.Fatalf("byte mismatch: %q != %q", c, min)
	}
	if err := r.UnreadByte(); err != nil {
		t.Fatal(err)
	}
	if err := r.UnreadByte(); err == nil {
		t.Fatal("expected UnreadByte error")
	}
}
