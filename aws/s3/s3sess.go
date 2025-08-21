package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Client struct {
	Key    string
	Secret string
	Region string
}

func (p *Client) S3Client() (s3api *s3.Client, err error) {
	s3api = s3.NewFromConfig(aws.Config{
		Region:      p.Region,
		Credentials: credentials.NewStaticCredentialsProvider(p.Key, p.Secret, ""),
	})

	return
}

func NewS3Client(params ...string) *Client {
	return &Client{
		Key:    params[0],
		Secret: params[1],
		Region: params[2],
	}
}
