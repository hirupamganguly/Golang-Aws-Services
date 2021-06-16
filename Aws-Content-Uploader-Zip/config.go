package awscontentuploader

import (
	"encoding/json"
	"os"
)

type Logging struct {
	FilePath string `json:"file_path"`
}

type Mongo struct {
	Credential *Credential `json:"credential"`
	Address    string      `json:"address"`
	DBName     string      `json:"db_name"`
}

type Credential struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Transport struct {
	HttpPort string `json:"http_port"`
}
type Config struct {
	Transport Transport `json:"transport"`
	Mongo     Mongo     `json:"mongo"`
	Logging   Logging   `json:"logging"`
	AWS       AWS       `json:"aws"`
	S3Bucket  S3Bucket  `json:"s3_bucket"`
}
type AWS struct {
	Region     string         `json:"region"`
	Credential *AWSCredential `json:"credential"`
	S3         *AWSS3         `json:"s3"`
}

type AWSCredential struct {
	CredentialType CredentialType `json:"credential_type"`
	ProfileName    string         `json:"profile_name"`
}
type CredentialType int

const (
	SharedConfig CredentialType = iota
	EC2IAM
)

type AWSS3 struct {
	ExerciseContentBucket           string `json:"exercise_content_bucket"`
	ContestExerciseBucketPrefixPath string `json:"exercisecontent_bucket_prefix_path"`
	PutReqExpirySignerInSec         int    `json:"put_req_expiry_signer_in_sec"`
}
type S3Bucket struct {
	Name string `json:"name"`
}

var DefaultConfig = Config{
	Transport: Transport{
		HttpPort: "3003",
	},
	AWS: AWS{
		Region: "ap-south-1",
		Credential: &AWSCredential{
			CredentialType: SharedConfig,
			ProfileName:    "RUPAM",
		},
		S3: &AWSS3{
			ExerciseContentBucket:           "",
			ContestExerciseBucketPrefixPath: "",
			PutReqExpirySignerInSec:         3600,
		},
	},
}

func Configure(cfgPath string) (c Config) {
	config, err := os.Open(cfgPath)
	if err != nil {
		return
	}
	defer config.Close()
	if err != nil {
		return
	}
	err = json.NewDecoder(config).Decode(&c)
	if err != nil {
		return
	}

	return
}

// r[0]
// strconv.ParseBool(r[1])
// r[2]
