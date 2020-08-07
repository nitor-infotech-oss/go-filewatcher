package main

import (
  "fmt"
  "github.com/nitor-infotech-oss/go-filewatcher/filewatcher"
  "github.com/nitor-infotech-oss/go-filewatcher/slacknotifications"
)

// Execution of filewatcher steps one after another
func main() {
  configObj := filewatcher.LoadConfiguration("config.json")
  log.Println("Fetching credentials")
  allFiles      := filewatcher.ListFiles(configObj)
  allFilesCount := len(allFiles)
  log.Println("Copying total ", allFilesCount, "files from ", configObj.SourceBucket," to ", configObj.DestBucket, " bucket")
  for i, filename := range allFiles {
    log.Println("---------------------------------------------------------")
    err := filewatcher.DecompressAndMove(configObj, filename)
    if (err == nil) {
      errArchival := filewatcher.ArchiveFile(configObj, filename)
      if (errArchival != nil) {
        log.Println("Unknown error occured.")
      }
    }
    log.Println("Moving ", i+1, " of ", allFilesCount, " files")
    log.Println("---------------------------------------------------------")
  }
  log.Println("All files copied and archived")
  msg := "Successfully moved files using filewatcher"
  err  = slacknotifications.SendSlackNotification(webhookUrl, msg)
  if err != nil {
    log.Println("Error while sending msg to slack")
  }
}
