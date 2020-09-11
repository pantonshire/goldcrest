package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "strings"
  "time"
)

func main() {
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

  //TODO: potential grouping of tweet get requests that come in at similar times?
  req, err := http.NewRequest(
    "GET",
    "https://api.twitter.com/2/tweets?ids=1261326399320715264,1278347468690915330",
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
