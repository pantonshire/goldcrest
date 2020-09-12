package proxy

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "path/filepath"
)

type Options struct {
  ConfigPath string
  Verbosity  int
}

type Config struct {
  GRPC struct {
    Enabled bool `json:"enabled"`
    Port    uint `json:"port"`
  } `json:"grpc"`
}

func Start(options Options) error {
  var configPath string
  if options.ConfigPath != "" {
    configPath = options.ConfigPath
  } else {
    defaultPaths := []string{
      "goldcrest.json",
      filepath.Join("conf", "goldcrest.json"),
    }
    for _, p := range defaultPaths {
      if exists, err := fileExists(p); err != nil {
        return err
      } else if exists {
        configPath = p
      }
    }
  }
  configData, err := ioutil.ReadFile(filepath.Clean(configPath))
  if err != nil {
    return err
  }
  var config Config
  if err := json.Unmarshal(configData, &config); err != nil {
    return err
  }
  //TODO: config is loaded now, use it elsewhere
  return nil
}

func fileExists(path string) (bool, error) {
  info, err := os.Stat(path)
  if os.IsNotExist(err) {
    return false, nil
  } else if err != nil {
    return false, err
  }
  return !info.IsDir(), nil
}
