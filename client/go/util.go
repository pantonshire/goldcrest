package goldcrest

import "strings"

func removeFromString(s string, cutAt ...Indices) string {
  n := uint(len(s))
  ignorePos := make([]bool, n)
  for _, indices := range cutAt {
    for i := indices.Start; i < indices.End && i < n; i++ {
      ignorePos[i] = true
    }
  }
  var buf strings.Builder
  for i, r := range s {
    if !ignorePos[i] {
      buf.WriteRune(r)
    }
  }
  return buf.String()
}
