package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// aws s3 putObject example: aws s3api put-object --bucket db-backup-huawen --key truco/rpmlist.txt(上传到bucket的位置) --body ./rpmlist.txt(本地文件)
func main() {

	sess, err := session.NewSession(&aws.Config{
		MaxRetries:  aws.Int(3),
		Credentials: credentials.NewSharedCredentials("C:/Users/Administrator/Desktop/aws.ini", "aws-huawen-root"),
		Region:      aws.String("sa-east-1"),
	})

	if err != nil {
		return
	}

	s3api := s3.New(sess)

	of, err := os.Open("C:/Users/Administrator/Desktop/baxi_MGCenter_FULL_20230211_000001.bak")
	if err != nil {
		return
	}

	defer of.Close()

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(of),
		Bucket: aws.String("db-backup-huawen"),
		Key:    aws.String("truco/baxi_MGCenter_FULL_20230211_000001.bak"),
	}

	res, err := s3api.PutObject(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)

}
