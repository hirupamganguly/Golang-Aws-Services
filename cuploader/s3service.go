package cuploader

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Storage struct {
	sess       *session.Session
	s3         *s3.S3
	bucket     string
	prefixPath string
}

func NewS3Storage(s *session.Session, bucket string, prefixPath string) *s3Storage {
	return &s3Storage{
		sess:       s,
		s3:         s3.New(s),
		bucket:     bucket,
		prefixPath: prefixPath,
	}
}

type Uploader interface {
	UploaderOfS3(path string, filePath string) (outPath string, version string, checksum string, modified bool, err error)
}

func (svc *s3Storage) UploaderOfS3(path string, filePath string) (outPath string, version string, checksum string, modified bool, err error) {
	rs, err := os.Open(filePath)
	if err != nil {
		return "", "", "", false, err
	}
	var key string
	if len(svc.prefixPath) > 0 {
		key = svc.prefixPath + "/" + path
	} else {
		key = path
	}

	hash := md5.New()
	_, err = io.Copy(hash, rs)
	if err != nil {
		return
	}

	checksum = hex.EncodeToString(hash.Sum(nil))
	reqHead, resHead := svc.s3.HeadObjectRequest(&s3.HeadObjectInput{
		Bucket: &svc.bucket,
		Key:    &key,
	})

	if err = reqHead.Send(); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case "Forbidden":
			case "NotFound":
			default:
				return "", "", "", false, err
			}
		} else {
			return "", "", "", false, err
		}

	} else {
		etag := strings.Trim(*resHead.ETag, "\"")
		if strings.Compare(etag, checksum) == 0 {
			if resHead.VersionId != nil {
				version = *resHead.VersionId
			}
			return reqHead.HTTPRequest.URL.String(), version, etag, false, nil
		}
	}

	_, err = rs.Seek(0, 0)
	if err != nil {
		return
	}

	req, res := svc.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &svc.bucket,
		Key:    &key,
		Body:   rs,
	})

	if err = req.Send(); err != nil {
		panic(err)
	}

	if strings.Compare(strings.Trim(*res.ETag, "\""), checksum) != 0 {
		svc.s3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: &svc.bucket,
			Key:    &key,
		})
		return "", "", "", false, err
	}

	if resHead.VersionId != nil {
		version = *resHead.VersionId
	}
	return req.HTTPRequest.URL.String(), version, checksum, true, nil
}

// r[2][:24]
// r[2][24:26]
// r[2][26:27]
// strconv.Atoi(r[2][27:28])
// r[2][28:30]
