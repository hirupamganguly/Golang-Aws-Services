package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3session *s3.S3

const (
	REGION      = "ap-south-1"
	BUCKET_NAME = "test.userrupam.bucket"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

type InputEvent struct {
	Link string `json:"link"`
	Key  string `json:"key"`
}

func GetImage(url string) (by []byte) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	by, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return by
}
func Handler(event InputEvent) (string, error) {
	image := GetImage(event.Link)
	_, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(image),
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(event.Key),
	})
	if err != nil {
		return "", err
	}
	return "Suc", err
}
func main() {
	lambda.Start(Handler)
}
