package utils

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(file io.Reader, filename string) (string, error) {
    // Convert io.Reader to io.ReadSeeker using a buffer
    buf := new(bytes.Buffer)
    if _, err := io.Copy(buf, file); err != nil {
        return "", err
    }
    fileReadSeeker := bytes.NewReader(buf.Bytes())

    s3session := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    }))
    s3svc := s3.New(s3session)

    bucket := "your-s3-bucket-name"
    key := fmt.Sprintf("uploads/%s", filename)

    _, err := s3svc.PutObject(&s3.PutObjectInput{
        Body:   fileReadSeeker,
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
        ACL:    aws.String("public-read"),
    })
    if err != nil {
        log.Println("Error uploading to S3:", err)
        return "", err
    }

    return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key), nil
}
