package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Object struct {
	S3Client *s3.Client
	Bucket   string
}

func (p *Object) PutLargeObject(src, dst string) (err error) {
	largeObject := []byte(src)
	largeBuffer := bytes.NewReader(largeObject)
	var partMiBs int64 = 10

	uploader := manager.NewUploader(p.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	input := &s3.PutObjectInput{
		Body:   largeBuffer,
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(dst),
	}

	_, err = uploader.Upload(context.TODO(), input)

	if err != nil {
		fmt.Println(err)
		return
	}

	return err
}
