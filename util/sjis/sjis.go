package sjis

// Small utility for Alpelo to encode/decode strings to/from Shift-JIS

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

var e *encoding.Encoder
var d *encoding.Decoder

type Lexeme struct {
	string string
	bytes  []byte
}

func init() {
	e = japanese.ShiftJIS.NewEncoder()
	d = japanese.ShiftJIS.NewDecoder()
}

func Encode(s string) []byte {
	t, err := e.String(s)
	if err != nil {
		return nil
	}
	return []byte(t)
}

func Decode(b []byte) string {
	t, err := d.Bytes(b)
	if err != nil {
		return ""
	}
	return string(t)
}

func NewString(s string) *Lexeme {
	return &Lexeme{
		s,
		Encode(s),
	}
}

func NewBytes(b []byte) *Lexeme {
	b = bytes.Trim(b, "\x00")
	return &Lexeme{
		Decode(b),
		b,
	}
}

func (s *Lexeme) Bytes() []byte {
	return append(s.bytes, 0)
}

func (s *Lexeme) String() string {
	return s.string
}

func (s *Lexeme) Length() int {
	return len(s.Bytes())
}
