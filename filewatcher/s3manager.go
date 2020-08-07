package filewatcher

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "fmt"
    "strings"
    // "bytes"
    "regexp"
)


// checks if the file stisfies regex conditions configured
func checkRegex(rules []string, filename string) bool {
  match := false
  for _, rule := range rules {
    match, _ := regexp.MatchString(rule, filename)
    if (match == true){
      return match
    }
  }
  return match
}

// get all the list of files from the s3 location
func ListFiles(configObj Config) []string {
  sess, err        := session.NewSession(&aws.Config{
      Region:      aws.String(configObj.Region),
      Credentials: credentials.NewStaticCredentials(configObj.AccessKey, configObj.AccessSecret, ""),
  })
  if err != nil {
    fmt.Println(err)
  }
  sourceBucket     := configObj.SourceBucket
  sourcePath       := configObj.SourcePath
  regexrules       := configObj.RegexRules

  svc := s3.New(sess)

  params := &s3.ListObjectsInput {
      Bucket: aws.String(sourceBucket),
      Prefix: aws.String(sourcePath),
  }

  allFilesArray := []string{}
  resp, _ := svc.ListObjects(params)
  for _, key := range resp.Contents {
      filename := strings.Replace(*key.Key, sourcePath+"/", "", 1)
      ismatching := checkRegex(regexrules, filename)
      if(ismatching == true) {
        allFilesArray = append(allFilesArray, filename)
      }
  }
  return allFilesArray
}
