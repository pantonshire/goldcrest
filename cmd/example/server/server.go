package main

import (
  "fmt"
  "goldcrest/proxy"
  "goldcrest/twitter1"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  ok, err, fatal, stop := proxy.ServeTwitter1(proxy.Twitter1Config{
    ProxyConfig: proxy.ProxyConfig{
      Enabled: true,
      Port:    7400,
    },
    TwitterConfig: twitter1.TwitterConfig{
      ClientTimeoutSeconds: 5,
    },
  })
  if err != nil {
    panic(err)
  } else if !ok {
    panic("Twitter1 server not started")
  }
  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(interrupt)
  select {
  case <-interrupt:
    fmt.Println("\nreceived interrupt, shutting down")
    stop()
    fmt.Println("goodbye!")
  case err := <-fatal:
    panic(err)
  }
}
