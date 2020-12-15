package proxy

// Returns the first non-empty string provided.
// If all strings are empty, an empty string is returned.
func strAlt(str string, strs ...string) string {
  if str != "" {
    return str
  }
  for _, str := range strs {
    if str != "" {
      return str
    }
  }
  return ""
}

func strSafeDeref(str *string) string {
  if str == nil {
    return ""
  }
  return *str
}
