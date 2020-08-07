package filewatcher

import (
	"log"
	"os"
  "strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "golang.org/x/sync/errgroup"
  "archive/zip"
  "path/filepath"
  "io"
)

type Downloader struct {
	manager           s3manager.Downloader
	bucket, key, dest string
}

// download file from S3 bucket
func NewDownloader(s *session.Session, bucket, key, dest string) *Downloader {
	return &Downloader{
		manager: *s3manager.NewDownloader(s),
		bucket:  bucket,
		key:     key,
		dest:    dest,
	}
}

func (d Downloader) Download() (string, error) {
	file, err := os.Create(d.dest)
	if err != nil {
		return "", err
	}
	defer file.Close()

	numBytes, err := d.manager.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    aws.String(d.key),
		})

	if err != nil {
		return "", err
	}
	log.Println("Unzipping - ", numBytes, "bytes")

	return file.Name(), nil
}

type Uploader struct {
	manager   s3manager.Uploader
	src, dest string
}

// Upload file to S3 bucket
func NewUploader(s *session.Session, src, dest string) *Uploader {
	return &Uploader{
		manager: *s3manager.NewUploader(s),
		src:     src,
		dest:    dest,
	}
}

func (u Uploader) Upload() error {
	eg := errgroup.Group{}

	err := filepath.Walk(u.src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		eg.Go(func() error {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			key := strings.Replace(file.Name(), u.src, "", 1)
			_, err = u.manager.Upload(&s3manager.UploadInput{
				Bucket: aws.String(u.dest),
				Key:    aws.String(key),
				Body:   file,
			})
			if err != nil {
				return err
			}
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Unzips file from one path and stores to another
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func DecompressAndMove(configObj Config, filename string) error {
	sess, err        := session.NewSession(&aws.Config{
      Region:      aws.String(configObj.Region),
      Credentials: credentials.NewStaticCredentials(configObj.AccessKey, configObj.AccessSecret, ""),
  })
  sourceBucket     := configObj.SourceBucket
  sourcePath       := configObj.SourcePath
  sourceFileName   := filename
  destBucket       := configObj.DestBucket
  sourceKey        := sourcePath + "/" + sourceFileName

  tempArtifactPath := "tmp/"
	tempZipPath      := tempArtifactPath + "zipped/"
	tempZip          := "temp.zip"
  tempUnzipPath    := tempArtifactPath + "unzipped/"

// download file and store at tmp location
  downloader := NewDownloader(sess, sourceBucket, sourceKey, tempZipPath+tempZip)
  downloadedZipPath, err := downloader.Download()
  if err != nil {
    log.Fatal(err)
  }

// unzip file and store at tmp location
  if err := Unzip(downloadedZipPath, tempUnzipPath); err != nil {
    log.Fatal(err)
  }

	// upload the unzipped file to destination bucket
  s3uploader := NewUploader(sess, tempUnzipPath, destBucket)

  if err := s3uploader.Upload(); err != nil {
    log.Fatal(err)
  }
	return nil
}
