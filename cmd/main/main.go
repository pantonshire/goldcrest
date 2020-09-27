package main

import (
  "encoding/json"
  "fmt"
  "github.com/jessevdk/go-flags"
  "goldcrest/proxy"
  "io/ioutil"
  "os"
  "os/signal"
  "path/filepath"
  "syscall"
)

var defaultConfigPaths = []string{"conf/goldcrest.json", "goldcrest.json"}

func main() {
  var opts struct {
    ConfigPath string `short:"c" long:"config" description:"Path to the config file to use"`
  }
  if _, err := flags.Parse(&opts); err != nil {
    if flagErr, ok := err.(*flags.Error); ok {
      if flagErr.Type == flags.ErrHelp {
        os.Exit(0)
      } else {
        os.Exit(1)
      }
    } else {
      panic(err)
    }
  }
  var possibleConfigPaths []string
  if opts.ConfigPath != "" {
    possibleConfigPaths = append(possibleConfigPaths, opts.ConfigPath)
  }
  possibleConfigPaths = append(possibleConfigPaths, defaultConfigPaths...)
  var configPath string
  for _, path := range possibleConfigPaths {
    cleaned := filepath.Clean(path)
    info, err := os.Stat(cleaned)
    if err != nil {
      if !os.IsNotExist(err) {
        panic(err)
      }
    } else if !info.IsDir() {
      configPath = cleaned
      break
    }
  }
  if configPath == "" {
    panic("could not find config file")
  }
  configData, err := ioutil.ReadFile(configPath)
  if err != nil {
    panic(err)
  }
  var config proxy.ServerConfig
  if err := json.Unmarshal(configData, &config); err != nil {
    panic(err)
  }
  ok, err, fatal, stop := proxy.StartTwitter1(config.Twitter1)
  if err != nil {
    panic(err)
  }
  if ok {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    defer signal.Stop(interrupt)
    select {
    case <-interrupt:
      stop()
    case err := <-fatal:
      panic(err)
    }
    fmt.Println("Goodbye!")
  }
}
