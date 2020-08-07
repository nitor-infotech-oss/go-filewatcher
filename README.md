# go-filewatcher

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

# Objective

- Develop a daemon based technology solution that will given a list of source watch directories and REGEX patterns for filenames.  Move files to a corresponding list of target directories with decompression, if the source file is compressed.
- This solution would in addition send alerts via slack to inform a group of recipients that a move event was triggered and the details of the event.  
- The list of source, target directories and rules should be configurable and stored in json file.

# Solution Approach
AWS Lambda is a serverless offering by AWS cloud and identified as s best fit for File Watcher components.
Below are the important components of the solution –
- This lambda function is the fileWatcher program which can be invoked using the api endpoint from external program.
- It reads the configurations from json file and perform actions on the file like unzip, move it to target bucket by creating the require directory structure on the target bucket.
- Logs the execution of the process in the and send an alert with execution details to Slack

### Dependencies

GoLang packages dependencies:
* [aws-sdk-go](https://github.com/aws/aws-sdk-go)
* [golang-x](https://golang.org/x/sync)

# Configurations

  - Settings.json - This will work as configuration file for golang package and contains below information.
    Credentials

        {
          "naming": {
              "prepend": "autogen"
          },
          "key": "XXXXXXXXXXXXXXXXXXXXXXXXXXX", <aws access key>
          "secret": "XXXXXXXXXXXXXXXXXXXXXXXX", <aws secret key>
          "region": "us-east-2", <aws region>
          "sourcebucket": "bucket-001", <source bucket name>
          "sourcepath": "inputpath", <source key path>
          "archivalpath": "archival", <path where file to be archived as is>
          "destbucket": "bucket-002", <destination bucket name>
          "destpath": "outputpath", <destination key path>
          "webhookurl": "https://hooks.slack.com/services/TKDC1KVV1/XXXXXX/XXXXXXXXXXXXXX", <slack webhook url>
          "regexrules": ["^20200801"] <slack webhook url>
        }}

### Imprting Package
```
import (
  "github.com/nitor-infotech-oss/go-filewatcher/filewatcher"
)
```

Package can be imported by simply using the path of repo where it resides.

### Steps to execute
- Setup $GOPATH
- Install required dependencies
```
    go run filewatcher.go
```

### Deploy
- Deployed over Amazon ec2 instance
- Can be deployed lambda function to make the executions faster
