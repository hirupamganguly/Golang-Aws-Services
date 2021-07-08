package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *s3.S3
)

const (
	REGION = "ap-south-1"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewSharedCredentials("", "tsm"),
	})))
}
func ListBucket() *s3.ListBucketsOutput {
	response, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	return response
}
func CreateBucket(name string) *s3.CreateBucketOutput {
	response, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION),
		},
	})
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok {
			switch awsError.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket Already exist")
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket Exist and Owned By Me")
			default:
				panic(err)
			}
		}
	}
	return response
}
func UploadObject(name string) *s3.PutObjectOutput {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	fmt.Println("uploading this file ", name, " ...")
	response, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String("test.rupam.bucket"),
		Key:    aws.String(strings.Split(name, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		panic(err)
	}
	return response
}
func ListobjectsInsideOfBucket() *s3.ListObjectsV2Output {
	response, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("test.rupam.bucket"),
	})
	if err != nil {
		panic(err)
	}
	return response
}
func downloadObjectFromBucket(name string) {
	fmt.Println("Downloading this file ", name, " ...")
	response, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("test.rupam.bucket"),
		Key:    aws.String(name),
	})
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	err = ioutil.WriteFile(name, body, 0644)
	if err != nil {
		panic(err)
	}
}
func DeleteObjectFromBucket(name string) *s3.DeleteObjectOutput {
	fmt.Println("Deleteing...")
	response, err := s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("test.rupam.bucket"),
		Key:    aws.String(name),
	})
	if err != nil {
		panic(err)
	}
	return response
}
func main() {
	// fmt.Println(CreateBucket("test.rupam.bucket"))
	// fmt.Println(ListBucket())
	// fmt.Println(UploadObject("assets/a1.png"))
	// fmt.Println(ListobjectsInsideOfBucket())
	// downloadObjectFromBucket("a1.png")
	// fmt.Println(DeleteObjectFromBucket("a1.png"))
}
