package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	config "repos/Golang-Aws-Services/Aws-Content-Uploader-Zip"
	"repos/Golang-Aws-Services/Aws-Content-Uploader-Zip/cuploader"
	contentuploader "repos/Golang-Aws-Services/Aws-Content-Uploader-Zip/upload"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

// /home/rupam/go/src/repos/Go-Lang-Backend-NoSQL-MySQL/AwsContentUploader/config.go

const cfgPath = "config.json"

func main() {
	var logger log.Logger
	ctx := context.Background()
	logger = log.NewJSONLogger(os.Stderr)
	httpLogger := log.With(logger, "component", "http")
	c := config.Configure(cfgPath)
	//fieldKeys := []string{"method"}
	awsSession, _ := session.NewSession(&aws.Config{
		Region:      aws.String(c.AWS.Region),
		Credentials: credentials.NewSharedCredentials("", c.AWS.Credential.ProfileName),
	})
	exerciseContentS3Store := cuploader.NewS3Storage(awsSession, c.AWS.S3.ExerciseContentBucket, c.AWS.S3.ContestExerciseBucketPrefixPath)
	awsservice := contentuploader.NewService(exerciseContentS3Store)
	mux := http.NewServeMux()
	mux.Handle("/contentuploader/", contentuploader.MakeHandler(ctx, awsservice, httpLogger))
	http.Handle("/", accessControl(mux))
	server := &http.Server{
		Addr:        ":" + c.Transport.HttpPort,
		IdleTimeout: time.Second * 30,
	}
	erChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		erChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "http", "address", ":"+c.Transport.HttpPort, "msg", "listening")

		erChan <- server.ListenAndServe()

	}()
	fmt.Println(<-erChan)
}
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
