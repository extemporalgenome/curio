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
//  * Printable bytes will be matched literally.
//  * \x00 will match zero or more whitespace bytes.
//  * \x01 will match zero or more non-whitespace bytes.
//  * Use of other non-printable or non-ASCII bytes is undefined.
//  * The behavior of printable pattern bytes immediately following \x01 is
//    undefined (currently, they will be captured in the \x01 group).
//  * The behavior of runs of identical capture bytes, such as `\x01\x01`,
//    is undefined.
func Match(s io.ByteScanner, p io.ByteReader, m []string) (n, o int, err error) {
	buf := make([]byte, 0, 1024)
	var b, c byte
	var v int
	for {
		b, err = p.ReadByte()
		o++
		switch {
		case err == io.EOF:
			err = nil
			fallthrough
		case err != nil:
			return
		}
		switch b {
		case '\x00':
			for {
				c, err = s.ReadByte()
				if err != nil {
					return
				}
				n++
				switch c {
				case ' ', '\n', '\t', '\r', '\f', '\v':
					continue
				}
				break
			}
			if err = s.UnreadByte(); err != nil {
				return
			}
			n--
		case '\x01':
			for {
				c, err = s.ReadByte()
				if err != nil {
					break
				}
				n++
				switch c {
				default:
					if m != nil {
						buf = append(buf, c)
					}
					continue
				case ' ', '\n', '\t', '\r', '\f', '\v':
				}
				if err = s.UnreadByte(); err == nil {
					n--
				}
				break
			}
			if m != nil {
				if _, ok := s.(*rb); ok {
					revbytes(buf)
				}
				m[v] = string(buf)
			}
			v++
			buf = buf[:0]
		default:
			c, err = s.ReadByte()
			if err != nil {
				return
			}
			n++
			if c == b {
				continue
			}
			if err = s.UnreadByte(); err == nil {
				n--
			}
			return
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
