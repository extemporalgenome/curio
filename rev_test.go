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

func TestRevByte(t *testing.T) {
	s := strings.NewReader("abcdefghijklmnopq")
	r := NewRevByteScanner(s, int64(s.Len()-1))
	const min, max byte = 'a', 'q'
	if err := r.UnreadByte(); err == nil {
		t.Fatal("expected UnreadByte error")
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

func ExampleNewRevByteScanner() {
	s := strings.NewReader("abcdefghijklmnopqrstuvwxyz")
	r := NewRevByteScanner(s, int64(s.Len()-1))
	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%c", c)
	}
	if err := r.UnreadByte(); err != nil {
		fmt.Println(err)
	} else if c, err := r.ReadByte(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%c", c)
	}
	// Output: zyxwvutsrqponmlkjihgfedcbaa
}

func ExampleNewRevRuneScanner() {
	s := strings.NewReader("你好世界")
	r := NewRevRuneScanner(s, int64(s.Len()-1))
	for {
		r, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%c", r)
	}
	if err := r.UnreadRune(); err != nil {
		fmt.Println(err)
	} else if r, _, err := r.ReadRune(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%c", r)
	}
	// Output: 界世好你你
}
