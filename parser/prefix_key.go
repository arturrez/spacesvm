// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package parser

import (
	"bytes"
	"errors"
)

const (
	MaxPrefixSize = 256
	MaxKeySize    = 256
	Delimiter     = byte(0x2f) // '/'
)

var (
	ErrPrefixEmpty = errors.New("prefix cannot be empty")
	ErrKeyEmpty    = errors.New("key cannot be empty")

	ErrPrefixTooBig = errors.New("prefix too big")
	ErrKeyTooBig    = errors.New("key too big")

	ErrInvalidDelimiter = errors.New("prefix/key has unexpected delimiters; only one sub-key is supported")
)

// CheckPrefix returns an error if the prefix format is invalid.
func CheckPrefix(pfx []byte) error {
	if len(pfx) == 0 {
		return ErrPrefixEmpty
	}
	if len(pfx) > MaxPrefixSize {
		return ErrPrefixTooBig
	}
	if bytes.Contains(pfx, []byte{Delimiter}) {
		return ErrInvalidDelimiter
	}
	return nil
}

// CheckKey returns an error if the key format is invalid.
func CheckKey(k []byte) error {
	if len(k) == 0 {
		return ErrKeyEmpty
	}
	if len(k) > MaxKeySize {
		return ErrKeyTooBig
	}
	if bytes.Contains(k, []byte{Delimiter}) {
		return ErrInvalidDelimiter
	}
	return nil
}

var noPrefixEnd = []byte{0}

// ParsePrefixKey parses the given string with delimiter to split prefix and key.
// "end" is the range end that can be used for the prefix query with "k".
func ParsePrefixKey(s []byte, opts ...OpOption) (pfx []byte, k []byte, end []byte, err error) {
	idx := bytes.IndexRune(s, rune(Delimiter))
	switch {
	case idx == -1: // "foo"
		pfx = s
	case idx == len(s)-1: // "foo/"
		pfx = s[:len(s)-1]
	default: // "a/b", then "a" becomes prefix, "b" becomes prefix
		splits := bytes.Split(s, []byte{Delimiter})
		pfx = splits[0]
		k = splits[1]
	}

	ret := &Op{}
	ret.applyOpts(opts)
	if ret.checkPrefix {
		if err = CheckPrefix(pfx); err != nil {
			return nil, nil, nil, err
		}
	}
	if ret.checkKey {
		if err = CheckKey(k); err != nil {
			return nil, nil, nil, err
		}
	}

	// next lexicographical key (range end) for prefix queries
	end = GetRangeEnd(k)
	return pfx, k, end, nil
}

// GetRangeEnd returns next lexicographical key (range end) for prefix queries
func GetRangeEnd(k []byte) (end []byte) {
	end = make([]byte, len(k))
	copy(end, k)
	pfxEndExist := false
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < 0xff {
			end[i]++
			end = end[:i+1]
			pfxEndExist = true
			break
		}
	}
	if !pfxEndExist {
		// next prefix does not exist (e.g., 0xffff);
		// default to special end key
		end = noPrefixEnd
	}
	return end
}

type Op struct {
	checkPrefix bool
	checkKey    bool
}

type OpOption func(*Op)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func WithCheckPrefix() OpOption {
	return func(op *Op) {
		op.checkPrefix = true
	}
}

func WithCheckKey() OpOption {
	return func(op *Op) {
		op.checkKey = true
	}
}