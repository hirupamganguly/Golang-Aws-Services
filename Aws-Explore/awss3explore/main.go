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
		Credentials: credentials.NewSharedCredentials("", "user_rupam"),
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
		Bucket: aws.String("test.userrupam.bucket"),
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
		Bucket: aws.String("test.userrupam.bucket"),
	})
	if err != nil {
		panic(err)
	}
	return response
}
func downloadObjectFromBucket(name string) {
	fmt.Println("Downloading this file ", name, " ...")
	response, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("test.userrupam.bucket"),
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
		Bucket: aws.String("test.userrupam.bucket"),
		Key:    aws.String(name),
	})
	if err != nil {
		panic(err)
	}
	return response
}
func main() {
	// fmt.Println(CreateBucket("test.userrupam.bucket"))
	// fmt.Println(ListBucket())
	// fmt.Println(UploadObject("assets/a1.png"))
	// fmt.Println(ListobjectsInsideOfBucket())
	// downloadObjectFromBucket("teacher.png")
	// fmt.Println(DeleteObjectFromBucket("a1.png"))
	// folder := "assets"
	// files, _ := ioutil.ReadDir(folder)
	// for _, file := range files {
	// 	if file.IsDir() {
	// 		continue
	// 	} else {
	// 		UploadObject(folder + "/" + file.Name())
	// 	}
	// }
	// fmt.Println(ListobjectsInsideOfBucket())
	// for _, object := range ListobjectsInsideOfBucket().Contents {
	// 	downloadObjectFromBucket(*object.Key)
	// 	// DeleteObjectFromBucket(*object.Key)
	// }
}

// OUTPUT:

// uploading this file  assets/a1.png  ...
// uploading this file  assets/s.png  ...
// uploading this file  assets/video.png  ...
// {
//   Contents: [{
//       ETag: "\"1ffcd585035f9325d66725ae7906ee8f\"",
//       Key: "a1.png",
//       LastModified: 2021-07-08 10:17:51 +0000 UTC,
//       Size: 177923,
//       StorageClass: "STANDARD"
//     },{
//       ETag: "\"14b3209a0bdb3f40dc4b00532064b88b\"",
//       Key: "s.png",
//       LastModified: 2021-07-08 10:17:52 +0000 UTC,
//       Size: 1985,
//       StorageClass: "STANDARD"
//     },{
//       ETag: "\"6d646ecd4eb652f5202d460902e30f79\"",
//       Key: "video.png",
//       LastModified: 2021-07-08 10:17:52 +0000 UTC,
//       Size: 11334,
//       StorageClass: "STANDARD"
//     }],
//   IsTruncated: false,
//   KeyCount: 3,
//   MaxKeys: 1000,
//   Name: "test.rupam.bucket",
//   Prefix: ""
// }
// Downloading this file  a1.png  ...
// Deleteing...
// Downloading this file  s.png  ...
// Deleteing...
// Downloading this file  video.png  ...
// Deleteing...
