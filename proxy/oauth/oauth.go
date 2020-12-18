package oauth

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

type Request struct {
  Method, Protocol, Domain, Path string
  Query, Body                    map[string]string
}

func NewRequest(method, protocol, domain, path string, query, body map[string]string) Request {
  return Request{
    Method:   method,
    Protocol: protocol,
    Domain:   domain,
    Path:     path,
    Query:    query,
    Body:     body,
  }
}

// Creates a new http.Request containing an authentication header as described at
// https://developer.twitter.com/en/docs/authentication/oauth-1-0a/authorizing-a-request
func (or Request) MakeRequest(ctx context.Context, secret, auth Auth) (*http.Request, error) {
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
