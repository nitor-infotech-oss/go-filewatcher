package filewatcher

import (
    "encoding/json"
    "os"
    "fmt"
)

// structs for holding config details
type Config struct {
  NamingRules struct {
      prepend     string `json:"prepend"`
  } `json:"naming"`
  AccessKey string `json:"key"`
  AccessSecret string `json:"secret"`
  Region         string `json:"region"`
  SourceBucket   string `json:"sourcebucket"`
  SourcePath     string `json:"sourcepath"`
  ArchivalPath   string `json:"archivalpath"`
  DestBucket     string `json:"destbucket"`
  DestPath       string `json:"destpath"`
  WebhookUrl     string `json:"webhookurl"`
  RegexRules    []string `json:"regexrules"`
}

// method for loading configurations
func LoadConfiguration(file string) Config {
    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}
