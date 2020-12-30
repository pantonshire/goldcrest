package main

import (
  "fmt"
  "github.com/jessevdk/go-flags"
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy"
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "net"
  "os"
  "os/signal"
  "syscall"
  "time"
)

const defaultConfigPath = "goldcrest.yaml"

type config struct {
  Server struct {
    Port           uint          `yaml:"port"`
    ConnectTimeout time.Duration `yaml:"connect_timeout"`
    TLS            struct {
      Enabled bool   `yaml:"enabled"`
      Crt     string `yaml:"crt"`
      Key     string `yaml:"key"`
    } `yaml:"tls"`
  } `yaml:"server"`
  Client struct {
    Timeout   time.Duration `yaml:"timeout"`
    Protocol  string        `yaml:"protocol"`
    BaseURL   string        `yaml:"base_url"`
    RateLimit struct {
      AssumeNext bool `yaml:"assume_next"`
    } `yaml:"rate_limit"`
  } `yaml:"client"`
}

func main() {
  var clfs struct {
    ConfigPath string `short:"c" long:"config" description:"Path to the config file to use"`
  }
  if _, err := flags.Parse(&clfs); err != nil {
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

  var configPath string
  if clfs.ConfigPath != "" {
    configPath = clfs.ConfigPath
  } else {
    configPath = defaultConfigPath
  }

  confData, err := ioutil.ReadFile(configPath)
  if err != nil {
    panic(err)
  }
  var conf config
  if err := yaml.Unmarshal(confData, &conf); err != nil {
    panic(err)
  }

  log := logrus.New()
  log.SetLevel(logrus.InfoLevel)

  address := fmt.Sprintf(":%d", conf.Server.Port)
  listener, err := net.Listen("tcp", address)
  if err != nil {
    panic(err)
  }

  log.Info("Listening at " + address)

  var opts []grpc.ServerOption

  if conf.Server.ConnectTimeout > 0 {
    opts = append(opts, grpc.ConnectionTimeout(conf.Server.ConnectTimeout))
  }

  if conf.Server.TLS.Enabled {
    creds, err := credentials.NewServerTLSFromFile(conf.Server.TLS.Crt, conf.Server.TLS.Key)
    if err != nil {
      panic(err)
    }
    opts = append(opts, grpc.Creds(creds))
  }

  server := grpc.NewServer(opts...)

  prox := proxy.NewProxy(
    log,
    conf.Client.Timeout,
    conf.Client.Protocol,
    conf.Client.BaseURL,
    conf.Client.RateLimit.AssumeNext,
  )
  pb.RegisterTwitterServer(server, prox)

  fatal := make(chan error, 1)

  go func() {
    if err := server.Serve(listener); err != nil {
      fatal <- err
    }
  }()

  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(interrupt)

  select {
  case <-interrupt:
    log.Info("Shutting down")
    server.GracefulStop()
    log.Info("Goodbye!")
  case err := <-fatal:
    panic(err)
  }
}
