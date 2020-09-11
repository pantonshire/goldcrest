package main

import (
  "bufio"
  "bytes"
  "crypto/hmac"
  "crypto/rand"
  "crypto/sha1"
  "encoding/base64"
  "fmt"
  "github.com/martinlindhe/base36"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "path"
  "sort"
  "strings"
  "time"
)

func main() {
  //v1()
  //v2()
  oauthParams := url.Values{}
  //These are example keys and tokens, not real!
  oauthParams.Set("oauth_consumer_key", "xvz1evFS4wEEPTGEFPHBog")
  oauthParams.Set("oauth_nonce", "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg")
  oauthParams.Set("oauth_signature_method", "HMAC-SHA1")
  oauthParams.Set("oauth_timestamp", "1318622958")
  oauthParams.Set("oauth_token", "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb")
  oauthParams.Set("oauth_version", "1.0")

  queryParams := url.Values{}
  queryParams.Set("include_entities", "true")

  bodyParams := url.Values{}
  bodyParams.Set("status", "Hello Ladies + Gentlemen, a signed OAuth request!")

  fmt.Println(getOAuthSignature(
    //Again, example key and token
    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
    "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
    "POST",
    "https://api.twitter.com/1.1/statuses/update.json",
    oauthParams,
    queryParams,
    bodyParams,
  ))
}

func v1() {
  client := http.Client{
    Timeout: time.Second * 5,
  }

  version := "1.1"
  method := "POST"
  reqPath := "statuses/update.json"
  baseURL := path.Join("https://api.twitter.com", version, reqPath)
  parameters := map[string]string{
    "status": "Hello, world!",
  }

  vals := url.Values{}
  for key, value := range parameters {
    vals.Set(key, value)
  }
  paramStr := vals.Encode() //TODO: use percentEncode

  reader := bufio.NewReader(os.Stdin)

  consumerKey, secretKey, token, tokenSecret :=
    readLn(reader, "consumer key"),
    readLn(reader, "secret key"),
    readLn(reader, "access token"),
    readLn(reader, "token secret")

  oauthVals := url.Values{}
  oauthVals.Set("oauth_consumer_key", consumerKey)
  oauthVals.Set("oauth_token", token)
  oauthVals.Set("oauth_signature_method", "HMAC-SHA1")
  oauthVals.Set("oauth_version", "1.0")
  oauthVals.Set("oauth_timestamp", fmt.Sprintf("%d", time.Now().Unix()))

  randBytes := make([]byte, 32)
  _, err := rand.Read(randBytes)
  if err != nil {
    panic(err)
  }
  oauthVals.Set("oauth_nonce", base36.EncodeBytes(randBytes))

  signature := getOAuthSignature(secretKey, tokenSecret, method, baseURL, oauthVals, nil, vals)
  oauthVals.Set("oauth_signature", signature)

  var dstBuilder strings.Builder
  dstBuilder.WriteString("OAuth ")
  {
    keys := make([]string, len(oauthVals))
    var i int
    for k := range oauthVals {
      keys[i] = k
      i++
    }
    sort.Strings(keys)
    var writtenFirst bool
    for _, key := range keys {
      vals := oauthVals[key]
      escapedKey := percentEncode(key)
      for _, val := range vals {
        if writtenFirst {
          dstBuilder.WriteString(", ")
        } else {
          writtenFirst = true
        }
        dstBuilder.WriteString(escapedKey)
        dstBuilder.WriteRune('=')
        dstBuilder.WriteRune('"')
        dstBuilder.WriteString(percentEncode(val))
        dstBuilder.WriteRune('"')
      }
    }
  }
  dst := dstBuilder.String()

  req, err := http.NewRequest(method, baseURL, bytes.NewBufferString(paramStr))
  if err != nil {
    panic(err)
  }

  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Set("Authorization", dst)

  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  for key, value := range resp.Header {
    fmt.Println(fmt.Sprintf("%s: %s", key, value))
  }

  respBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }

  fmt.Println(string(respBody))
}

func v2() {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("bearer token :> ")
  bearerToken, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  bearerToken = strings.TrimSpace(bearerToken)

  client := http.Client{
    Timeout: time.Second * 5,
  }

  req, err := http.NewRequest(
    "GET",
    //"https://api.twitter.com/2/tweets?ids=1261326399320715264,1278347468690915330",
    "https://api.twitter.com/2/tweets?ids=1304387795507650561&tweet.fields=created_at,conversation_id,attachments,entities&expansions=author_id,attachments.media_keys",
    nil,
  )
  if err != nil {
    panic(err)
  }
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  for key, value := range resp.Header {
    fmt.Println(fmt.Sprintf("%s: %s", key, value))
  }

  respBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }

  fmt.Println(string(respBody))
}

func getOAuthSignature(consumerSecret, tokenSecret, method, baseURL string, oauthParams, queryParams, bodyParams map[string][]string) string {
  allParams := make(map[string][]string)
  for _, params := range []url.Values{oauthParams, queryParams, bodyParams} {
    for key, values := range params {
      for _, value := range values {
        allParams[key] = append(allParams[key], value)
      }
    }
  }

  var paramBuilder strings.Builder
  {
    keys := make([]string, len(allParams))
    var i int
    for k := range allParams {
      keys[i] = k
      i++
    }
    sort.Strings(keys)
    var writtenFirst bool
    for _, key := range keys {
      vals := allParams[key]
      escapedKey := percentEncode(key)
      for _, val := range vals {
        if writtenFirst {
          paramBuilder.WriteString("&")
        } else {
          writtenFirst = true
        }
        paramBuilder.WriteString(escapedKey)
        paramBuilder.WriteRune('=')
        paramBuilder.WriteString(percentEncode(val))
      }
    }
  }
  paramStr := paramBuilder.String()

  sigBase := strings.ToUpper(method) + "&" + percentEncode(baseURL) + "&" + percentEncode(paramStr)
  signingKey := percentEncode(consumerSecret) + "&" + percentEncode(tokenSecret)

  h := hmac.New(sha1.New, []byte(signingKey))
  h.Write([]byte(sigBase))

  return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func readLn(reader *bufio.Reader, prompt string) string {
  fmt.Print(fmt.Sprintf("%s :> ", prompt))
  str, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  return strings.TrimSpace(str)
}

func percentEncode(raw string) string {
  var dst []byte
  for _, b := range []byte(raw) {
    if skipPercentageEncode(b) {
      dst = append(dst, b)
    } else {
      b1, b2 := toHex(b)
      dst = append(dst, 0x25, b1, b2)
    }
  }
  return string(dst)
}

func skipPercentageEncode(r byte) bool {
  return (0x30 <= r && r <= 0x39) || (0x41 <= r && r <= 0x5A) || (0x61 <= r && r <= 0x7A) ||
    r == 0x2D || r == 0x2E || r == 0x5F || r == 0x7E
}

func toHex(b byte) (b1, b2 byte) {
  if b1 = (b >> 4) + 0x30; b1 > 0x39 {
    b1 += 0x7
  }
  if b2 = (b & 0xF) + 0x30; b2 > 0x39 {
    b2 += 0x7
  }
  return b1, b2
}
