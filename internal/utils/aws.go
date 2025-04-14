package utils

import (
	"bytes"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/aws/aws-sdk-go/service/s3"
)


func S3PutObject(fileEncoded []byte, bucketName *string, keyName *string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		panic(err)
	}

	s3Uploader := s3manager.NewUploader(sess)

	// Upload input parameters
	uploadParams := &s3manager.UploadInput{
		Bucket: aws.String(*bucketName),
		Key:    aws.String(*keyName),
		Body:   bytes.NewReader(fileEncoded),
	}

	// Perform an upload.
	_, err = s3Uploader.Upload(uploadParams)
	if err != nil {
		panic(err)
	}

}

func S3GetPresignedURL(bucketName *string, keyName *string) *string {

    sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		panic(err)
	}

    // Create S3 service client
    s3Clent := s3.New(sess)

	requestParams := &s3.GetObjectInput{
        Bucket: aws.String(*bucketName),
		Key:    aws.String(*keyName),
    }

    req, _ := s3Clent.GetObjectRequest(requestParams)
    urlStr, err := req.Presign(5 * time.Minute)
    if err != nil {
        panic(err)
    }

    return &urlStr
}
