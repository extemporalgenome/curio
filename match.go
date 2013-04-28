// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curio

import "io"

// Match provides a simple sequential ASCII text matcher.
// It is specialized for processing well formed, structured printable ASCII.
// s yields input bytes and p yields pattern bytes.
// n and o correspond to the number of bytes read from s and p, respectively.
// If m is non-nil, len(m) should be greater than or equal to the number of
// capture-bytes in the pattern.
//
// Patterns:
//  * Printable ASCII bytes will be matched literally.
//  * All groups (specified by non-printable bytes in the pattern stream)
//    are non-greedy, and match zero or more characters.
//  * \x00: ASCII whitespace bytes.
//  * \x01: non-whitespace ASCII bytes.
//  * \xd0: decimal digits.
//  * \xd6: hexadecimal digits.
//  * \xd8: octal digits.
//  * \xd3: base-36 digits.
//  * \xfe: printable ASCII bytes.
//  * \xff: 8-bit bytes.
//  * All groups are capturing beside \x00.
//  * Use of other non-printable or non-ASCII bytes is undefined.
func Match(s, p io.ByteScanner, m []string) (n, o int, err error) {
	var (
		quit = false
		buf  = make([]byte, 0, 1024)
		a, b byte
		c    byte
		v    int
	)
	for !quit {
		a, err = p.ReadByte()
		switch {
		case err == io.EOF:
			err = nil
			fallthrough
		case err != nil:
			return
		}
		o++
		if tab[a]&prg == 0 {
			c, err = s.ReadByte()
			if err != nil {
				return
			}
			n++
			if a == c {
				continue
			} else if err = s.UnreadByte(); err == nil {
				err = ErrByteMismatch
				n--
			}
			return
		}
		a = tab[a]
		b, err = p.ReadByte()
		if err == io.EOF {
			b = nop
		} else if err != nil {
			return
		} else if err = p.UnreadByte(); err != nil {
			o++
			return
		}
		for {
			c, err = s.ReadByte()
			if err != nil {
				quit = true
				goto fill
			}
			n++
			if tab[b]&prg == 0 && c == b ||
				tab[b]&prg != 0 && tab[c]&tab[b] != 0 ||
				a != any && tab[c]&a == 0 {
				break
			} else if a&^prg != ws {
				buf = append(buf, c)
			}
		}
		if err = s.UnreadByte(); err == nil {
			err = io.EOF
			n--
		}
	fill:
		if a&^prg != ws {
			if _, ok := s.(*rb); ok {
				revbytes(buf)
			}
			m[v] = string(buf)
			v++
			buf = buf[:0]
		}
	}
	return
}

func revbytes(p []byte) {
	j := len(p)
	m := j / 2
	for i := 0; i < m; i++ {
		j--
		p[i], p[j] = p[j], p[i]
	}
}

const (
	prg = 1 << iota
	ws
	pc
	ll
	lu
	do
	dd
	dx
	lx  = ll | lu
	wrd = lx | dd
	nws = wrd | pc
	prn = nws | ws
	any = 0xff
	nop = 0x80
)

var tab = [1 << 8]byte{
	' ': ws, '\n': ws, '\t': ws, '\r': ws, '\f': ws, '\v': ws,
	'0': do | dd | dx, '1': do | dd | dx, '2': do | dd | dx, '3': do | dd | dx,
	'4': do | dd | dx, '5': do | dd | dx, '6': do | dd | dx, '7': do | dd | dx,
	'8': dd | dx, '9': dd | dx,

	'a': ll | dx, 'b': ll | dx, 'c': ll | dx,
	'd': ll | dx, 'e': ll | dx, 'f': ll | dx,
	'g': ll, 'h': ll, 'i': ll, 'j': ll, 'k': ll,
	'l': ll, 'm': ll, 'n': ll, 'o': ll, 'p': ll,
	'q': ll, 'r': ll, 's': ll, 't': ll, 'u': ll,
	'v': ll, 'w': ll, 'x': ll, 'y': ll, 'z': ll,

	'A': lu | dx, 'B': lu | dx, 'C': lu | dx,
	'D': lu | dx, 'E': lu | dx, 'F': lu | dx,
	'G': lu, 'H': lu, 'I': lu, 'J': lu, 'K': lu,
	'L': lu, 'M': lu, 'N': lu, 'O': lu, 'P': lu,
	'Q': lu, 'R': lu, 'S': lu, 'T': lu, 'U': lu,
	'V': lu, 'W': lu, 'X': lu, 'Y': lu, 'Z': lu,

	'!': pc, '"': pc, '#': pc, '$': pc, '%': pc, '&': pc, '\'': pc, '(': pc,
	')': pc, '*': pc, '+': pc, ',': pc, '-': pc, '.': pc, '/': pc, ':': pc,
	';': pc, '<': pc, '=': pc, '>': pc, '?': pc, '@': pc, '[': pc, '\\': pc,
	']': pc, '^': pc, '_': pc, '`': pc, '{': pc, '|': pc, '}': pc, '~': pc,

	'\x00': prg | ws,
	'\x01': prg | nws,
	'\xd8': prg | do,
	'\xd0': prg | dd,
	'\xd6': prg | dx,
	'\xd3': prg | dd | lx,
	'\xfe': prg | prn,
	'\xff': any,
}
