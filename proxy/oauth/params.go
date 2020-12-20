package oauth

import (
  "sort"
  "strings"
)

type Params map[string]string
type percentEncodedParams Params

func NewParams() Params {
  return make(Params)
}

func (ps Params) Set(key, val string) bool {
  if _, ok := ps[key]; ok {
    return false
  }
  ps[key] = val
  return true
}

func (ps Params) Extend(other Params) {
  for k, v := range other {
    ps.Set(k, v)
  }
}

func (pp percentEncodedParams) encode(sep string, quote bool) string {
  n := len(pp)
  if n == 0 {
    return ""
  }
  keys := make([]string, n)
  var i int
  for k := range pp {
    keys[i] = k
    i++
  }
  sort.Strings(keys)
  var dst strings.Builder
  dst.WriteString(percentEncode(keys[0]))
  dst.WriteRune('=')
  if quote {
    dst.WriteRune('"')
    dst.WriteString(percentEncode(pp[keys[0]]))
    dst.WriteRune('"')
  } else {
    dst.WriteString(percentEncode(pp[keys[0]]))
  }
  for j := 1; j < n; j++ {
    dst.WriteString(sep)
    dst.WriteString(percentEncode(keys[j]))
    dst.WriteRune('=')
    if quote {
      dst.WriteRune('"')
      dst.WriteString(percentEncode(pp[keys[j]]))
      dst.WriteRune('"')
    } else {
      dst.WriteString(percentEncode(pp[keys[j]]))
    }
  }
  return dst.String()
}

// Encodes the given string according to RFC 3986, Section 2.1.
func percentEncode(raw string) string {
  var dst []byte
  for _, b := range []byte(raw) {
    if allowRawByte(b) {
      dst = append(dst, b)
    } else {
      b1, b2 := toHex(b)
      dst = append(dst, 0x25, b1, b2)
    }
  }
  return string(dst)
}

func allowRawByte(b byte) bool {
  return (0x30 <= b && b <= 0x39) || (0x41 <= b && b <= 0x5A) || (0x61 <= b && b <= 0x7A) ||
    b == 0x2D || b == 0x2E || b == 0x5F || b == 0x7E
}

func toHex(b byte) (b1, b2 byte) {
  if b1 = (b >> 4) + 0x30; b1 > 0x39 {
    b1 += 0x7
  }
  if b2 = (b & 0x0F) + 0x30; b2 > 0x39 {
    b2 += 0x7
  }
  return b1, b2
}
