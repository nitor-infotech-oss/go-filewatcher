package filewatcher

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "fmt"
    "os"
)

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}


// copy a file in as is state from one bucket to andestBucket
func ArchiveFile(configObj Config, filename string) error {
  sess, err        := session.NewSession(&aws.Config{
      Region:      aws.String(configObj.Region),
      Credentials: credentials.NewStaticCredentials(configObj.AccessKey, configObj.AccessSecret, ""),
  })
  sourceBucket     := configObj.SourceBucket
  sourcePath       := configObj.SourcePath
  sourceFileName   := filename
  destBucket       := configObj.DestBucket
  destPath         := configObj.DestPath
  sourceKey        := sourcePath + "/" + sourceFileName
  source           := sourceBucket + "/" + sourcePath + "/" + sourceFileName
  destKey          := destPath   + "/" + sourceFileName

  svc := s3.New(sess)

  _, err = svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(destBucket), CopySource: aws.String(source), Key: aws.String(destKey)})
  if err != nil {
      exitErrorf("Unable to copy item from bucket %q to bucket %q, %v", sourceBucket, destBucket, err)
  }

  // Wait to see if the item got copied
  err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(destBucket), Key: aws.String(sourceKey)})
  if err != nil {
      exitErrorf("Error occurred while waiting for item %q to be copied to bucket %q, %v", sourceBucket, sourceKey, destBucket, err)
  }
  return nil
}
