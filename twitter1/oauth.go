package twitter1

import (
  "bytes"
  "context"
  "crypto/hmac"
  "crypto/rand"
  "crypto/sha1"
  "encoding/base64"
  "fmt"
  "github.com/martinlindhe/base36"
  "net/http"
  "path"
  "sort"
  "strings"
  "time"
)

const (
  oauthVersion         = "1.0"
  oauthSignatureMethod = "HMAC-SHA1"
  oauthNonceBytes      = 32
)

type Auth struct {
  Key   string `json:"key"`
  Token string `json:"token"`
}

type AuthPair struct {
  Secret Auth `json:"secret"`
  Public Auth `json:"public"`
}

type OAuthRequest struct {
  Method, Protocol, Domain, Path string
  Query, Body                    map[string]string
}

type PercentEncodedParams map[string]string

// Creates a new http.Request containing an authentication header as described at
// https://developer.twitter.com/en/docs/authentication/oauth-1-0a/authorizing-a-request
func (or OAuthRequest) MakeRequest(ctx context.Context, secret, auth Auth) (*http.Request, error) {
  nonce, err := randBase36(oauthNonceBytes)
  if err != nil {
    return nil, err
  }

  baseURL := or.Protocol + "://" + path.Join(or.Domain, or.Path)

  queryParams, bodyParams := PercentEncodedParams(or.Query), PercentEncodedParams(or.Body)

  timestamp := fmt.Sprintf("%d", time.Now().Unix())

  oauthParams := PercentEncodedParams{}
  oauthParams.Set("oauth_consumer_key", auth.Key)
  oauthParams.Set("oauth_token", auth.Token)
  oauthParams.Set("oauth_signature_method", oauthSignatureMethod)
  oauthParams.Set("oauth_version", oauthVersion)
  oauthParams.Set("oauth_timestamp", timestamp)
  oauthParams.Set("oauth_nonce", nonce)

  signature := signOAuth(secret, or.Method, baseURL, oauthParams, queryParams, bodyParams)
  oauthParams.Set("oauth_signature", signature)

  authorization := "OAuth " + oauthParams.Encode(", ", true)

  fullURL := baseURL + "?" + queryParams.Encode("&", false)
  bodyStr := bodyParams.Encode("&", false)

  req, err := http.NewRequestWithContext(ctx, or.Method, fullURL, bytes.NewBufferString(bodyStr))
  if err != nil {
    return nil, err
  }

  if bodyStr != "" {
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  }
  req.Header.Set("Authorization", authorization)

  return req, nil
}

// Creates an OAuth signature using the method described at
// https://developer.twitter.com/en/docs/authentication/oauth-1-0a/creating-a-signature
func signOAuth(secret Auth, method, baseURL string, oauthParams, queryParams, bodyParams PercentEncodedParams) string {
  allParams := PercentEncodedParams{}
  for key, value := range oauthParams {
    allParams.Set(key, value)
  }
  for key, value := range queryParams {
    allParams.Set(key, value)
  }
  for key, value := range bodyParams {
    allParams.Set(key, value)
  }
  paramStr := allParams.Encode("&", false)
  sigBase := strings.ToUpper(method) + "&" + PercentEncode(baseURL) + "&" + PercentEncode(paramStr)
  sigKey := PercentEncode(secret.Key) + "&" + PercentEncode(secret.Token)
  hash := hmac.New(sha1.New, []byte(sigKey))
  hash.Write([]byte(sigBase))
  return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func randBase36(bytes int) (string, error) {
  randBytes := make([]byte, bytes)
  if _, err := rand.Read(randBytes); err != nil {
    return "", err
  }
  return base36.EncodeBytes(randBytes), nil
}

func (pp PercentEncodedParams) Set(key, val string) bool {
  if _, ok := pp[key]; ok {
    return false
  }
  pp[key] = val
  return true
}

func (pp PercentEncodedParams) Encode(sep string, quote bool) string {
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
  dst.WriteString(PercentEncode(keys[0]))
  dst.WriteRune('=')
  if quote {
    dst.WriteRune('"')
    dst.WriteString(PercentEncode(pp[keys[0]]))
    dst.WriteRune('"')
  } else {
    dst.WriteString(PercentEncode(pp[keys[0]]))
  }
  for j := 1; j < n; j++ {
    dst.WriteString(sep)
    dst.WriteString(PercentEncode(keys[j]))
    dst.WriteRune('=')
    if quote {
      dst.WriteRune('"')
      dst.WriteString(PercentEncode(pp[keys[j]]))
      dst.WriteRune('"')
    } else {
      dst.WriteString(PercentEncode(pp[keys[j]]))
    }
  }
  return dst.String()
}

// Encodes the given string according to RFC 3986, Section 2.1.
func PercentEncode(raw string) string {
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
