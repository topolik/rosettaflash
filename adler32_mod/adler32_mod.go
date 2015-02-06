/*
 * Copyright 2014 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package adler32_mod

import (
	"sort"
	"unicode/utf8"

	"github.com/mikispag/rosettaflash/charset"
)

const (
	// mod is the largest prime that is less than 65536.
	mod = 65521
	// nmax is the largest n such that
	// 255 * n * (n+1) / 2 + (n+1) * (mod-1) <= 2^32-1.
	// It is mentioned in RFC 1950 (search for "5552").
	nmax = 5552
)

func Checksum_allowed(checksum uint32, allowed_charset *charset.Charset) bool {
	S1, S2 := S1(checksum), S2(checksum)
	return S1_S2_allowed(S1, S2, allowed_charset)
}

func S1_S2_allowed(S1, S2 int, allowed_charset *charset.Charset) bool {
	return S_allowed_UTF8(S2) && S_allowed(S1, allowed_charset)
}

func S_allowed(S int, allowed_charset *charset.Charset) bool {
	combinations := (*allowed_charset).Combinations
	index_S := sort.SearchInts(combinations, S)

	if index_S == len(combinations) {
		return false
	}

	if combinations[index_S] == S {
		return true
	}
	return false
}

func S_allowed_UTF8(S int) bool {
	return S > 0xc080 && IsUTF8(S)
}

func IsUTF8(S int) bool {
	firstByte := (S & 0xff00) >> 8
	secondByte := S & 0xff

	utf8Sequence := []byte{byte(firstByte), byte(secondByte)}

	return utf8.Valid(utf8Sequence)
}

// Add p to the running checksum d.
func Update(d uint32, p []byte) uint32 {
	s1, s2 := uint32(d&0xffff), uint32(d>>16)
	for len(p) > 0 {
		var q []byte
		if len(p) > nmax {
			p, q = p[:nmax], p[nmax:]
		}
		for _, x := range p {
			s1 += uint32(x)
			s2 += s1
		}
		s1 %= mod
		s2 %= mod
		p = q
	}
	return uint32(s2<<16 | s1)
}

func S1(d uint32) int {
	return int(d & 0xffff)
}

func S2(d uint32) int {
	return int(d >> 16)
}

// Checksum returns the Adler-32 checksum of data.
func Checksum(data []byte) uint32 { return uint32(Update(1, data)) }
