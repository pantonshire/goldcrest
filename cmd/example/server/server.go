package main

import (
  "fmt"
  "github.com/pantonshire/goldcrest/proxy"
  "google.golang.org/grpc"
  "net"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  err, fatal, stop := serve(7400)
  if err != nil {
    panic(err)
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

func serve(port uint) (startupErr error, fatal <-chan error, shutdown func()) {
  var opts []grpc.ServerOption
  listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
  if err != nil {
    return err, nil, nil
  }
  grpcServer := grpc.NewServer(opts...)
  proxy.InitServer(grpcServer)
  fatalErrors := make(chan error, 1)
  go func() {
    if err := grpcServer.Serve(listener); err != nil {
      fatalErrors <- err
    }
  }()
  return nil, fatalErrors, func() {
    grpcServer.GracefulStop()
  }
}
