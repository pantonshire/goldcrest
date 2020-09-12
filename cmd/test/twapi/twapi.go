package main

import (
  "bufio"
  "fmt"
  "goldcrest/twitter1"
  "io/ioutil"
  "net/http"
  "os"
  "path"
  "strings"
  "time"
)

func main() {
  v1()
  //v2()
  //oauthParams := twitter1.PercentEncodedParams{}
  ////These are example keys and tokens, not real!
  //oauthParams.Set("oauth_consumer_key", "xvz1evFS4wEEPTGEFPHBog")
  //oauthParams.Set("oauth_nonce", "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg")
  //oauthParams.Set("oauth_signature_method", "HMAC-SHA1")
  //oauthParams.Set("oauth_timestamp", "1318622958")
  //oauthParams.Set("oauth_token", "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb")
  //oauthParams.Set("oauth_version", "1.0")
  //
  //queryParams := twitter1.PercentEncodedParams{}
  //queryParams.Set("include_entities", "true")
  //
  //bodyParams := twitter1.PercentEncodedParams{}
  //bodyParams.Set("status", "Hello Ladies + Gentlemen, a signed OAuth request!")
  //
  //fmt.Println(getOAuthSignature(
  //  //Again, example key and token
  //  "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
  //  "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
  //  "POST",
  //  "https://api.twitter.com/1.1/statuses/update.json",
  //  oauthParams,
  //  queryParams,
  //  bodyParams,
  //))
}

func v1() {
  client := http.Client{
    Timeout: time.Second * 5,
  }

  reader := bufio.NewReader(os.Stdin)

  consumerKey, secretKey, token, tokenSecret :=
    readLn(reader, "consumer key"),
    readLn(reader, "secret key"),
    readLn(reader, "access token"),
    readLn(reader, "token secret")

  req, err := twitter1.OAuthRequest{
    Method:   "POST",
    Protocol: "https",
    Domain:   "api.twitter.com",
    Path:     path.Join("1.1", "statuses/update.json"),
    Body: map[string]string{
      "status": "One final hello world, probably 😎",
    },
  }.MakeRequest(twitter1.Auth{Key: secretKey, Token: tokenSecret}, twitter1.Auth{Key: consumerKey, Token: token})
  if err != nil {
    panic(err)
  }

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

func readLn(reader *bufio.Reader, prompt string) string {
  fmt.Print(fmt.Sprintf("%s :> ", prompt))
  str, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  return strings.TrimSpace(str)
}
