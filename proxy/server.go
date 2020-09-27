package proxy

import (
  "fmt"
  "goldcrest/twitter1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "net"
)

type ServerConfig struct {
  Twitter1 Twitter1Config `json:"twitter1"`
}

type ProxyConfig struct {
  Enabled bool      `json:"enabled"`
  Port    uint      `json:"port"`
  TLS     TLSConfig `json:"tls"`
}

type TLSConfig struct {
  Enabled bool   `json:"enabled"`
  Cert    string `json:"cert"`
  Key     string `json:"key"`
}

type Twitter1Config struct {
  ProxyConfig
  twitter1.TwitterConfig
}

func ServeTwitter1(conf Twitter1Config) (ok bool, startupErr error, fatal <-chan error, shutdown func()) {
  if !conf.Enabled {
    return false, nil, nil, nil
  }
  var opts []grpc.ServerOption
  tlsOpts, err := conf.TLS.parse()
  if err != nil {
    return false, err, nil, nil
  }
  opts = append(opts, tlsOpts...)
  listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
  if err != nil {
    return false, err, nil, nil
  }
  grpcServer := grpc.NewServer()
  twitter := twitter1.NewTwitter(twitter1.TwitterConfig{ClientTimeoutSeconds: 5})
  if err := twitter.Server(grpcServer); err != nil {
    return false, err, nil, nil
  }
  fatalErrors := make(chan error, 1)
  go func() {
    if err := grpcServer.Serve(listener); err != nil {
      fatalErrors <- err
    }
  }()
  return true, nil, fatalErrors, func() {
    grpcServer.GracefulStop()
  }
}

func (conf TLSConfig) parse() ([]grpc.ServerOption, error) {
  if !conf.Enabled {
    return nil, nil
  }
  creds, err := credentials.NewServerTLSFromFile(conf.Cert, conf.Key)
  if err != nil {
    return nil, err
  }
  return []grpc.ServerOption{grpc.Creds(creds)}, nil
}
