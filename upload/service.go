package contentuploader

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"repos/AwsContentUploader/cuploader"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type UploadSvc interface {
	Upload(ctx context.Context, fileName string, r io.ReaderAt, size int64)
}
type Service interface {
	UploadSvc
}

type service struct {
	cup cuploader.Uploader
}

func NewService(cup cuploader.Uploader) *service {
	return &service{
		cup: cup,
	}
}
func unzip(rat io.ReaderAt, size int64, dest string) error {
	r, err := zip.NewReader(rat, size)
	if err != nil {
		return err
	}
	os.MkdirAll(dest, 0755)
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}
func extractInfoFromTextFile(fpath string) (extractInfos []ExtractInfo, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		if len(data) != 3 {
			continue
		}

		if len(data[2]) != 30 {
			continue
		}
		name := data[0]
		isTrue, _ := strconv.ParseBool(data[])
		payload := data[]
		dataId := data[][:]
		version := data[][:]
		itype := data[][:]
		gender, _ := strconv.Atoi(data[][:])
		label := data[][:]


		extractInfo := ExtractInfo{
			Name:    name,
			IsTrue:  isTrue,
			Payload: payload,
			DataId:  dataId,
			Version: version,
			Type:    itype,
			Gender:  Gender(gender),
			Label:   label,
		}
		extractInfos = append(extractInfos, extractInfo)
	}
	return extractInfos, nil
}
func (svc *service) Upload(ctx context.Context, fileName string, r io.ReaderAt, size int64) {
	var ExtractInfos []ExtractInfo
	parts := strings.Split(fileName[:len(fileName)-len(filepath.Ext(fileName))], "_")
	if len(parts) < 2 {
		fmt.Println("Does Not Contain '_")
	}
	tempDir := filepath.Join("/home/rupam/go/src/", bson.NewObjectId().Hex())
	err := unzip(r, size, tempDir)
	if err != nil {
		return
	}
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fName := filepath.Base(path)
		if ext := filepath.Ext(fName); ext == ".txt" {
			ExtractInfos, _ = extractInfoFromTextFile(path)
		}
		return nil
	})
	fmt.Println(ExtractInfos)
	// Upload data to AWS s3 and Get DataFile:-
	for _, extractInfo := range ExtractInfos {
		fileMediaPath, fileMetaPath := GetMediaAndMetaFilePath("", 0, extractInfo.Version, extractInfo.Gender, tempDir, extractInfo.Payload)
		dbMediaPath, dbMetaPath := GetDbMediaPath(extractInfo.Version, extractInfo.Gender, extractInfo.DataId, filepath.Base(fileMediaPath), filepath.Base(fileMetaPath))
		fMediaPath := strings.TrimSpace(fileMediaPath)
		fMetaPath := strings.TrimSpace(fileMediaPath)
		if len(fMediaPath) < 1 || len(fMetaPath) < 1 {
			return
		}
		S3Path, S3VerionId, Md5, ok, Err := svc.cup.UploaderOfS3(dbMediaPath, fMediaPath)
		fmt.Println(S3Path, "-> ", S3VerionId, "-> ", Md5, "-> ", ok, "-> ", Err, " ") // Add to DB
		S3Path, S3VerionId, Md5, ok, Err = svc.cup.UploaderOfS3(dbMetaPath, fMetaPath)
		fmt.Println(S3Path, "-> ", S3VerionId, "-> ", Md5, "-> ", ok, "-> ", Err, " ") // Add to DB
	}
	defer os.Remove(tempDir)
}

func GetMediaAndMetaFilePath(refId string, bitType int, version string, gender Gender, tempDir string, payload string) (mediaPath string, metaPath string) {
	return fmt.Sprintf("%s/d_%s", tempDir, payload), fmt.Sprintf("%s/m_%s", tempDir, payload)
}

func GetDbMediaPath(version string, gender Gender, schoolId string, fMediaPath string, fMetaPath string) (fileMediaPath string, fileMetaPath string) {
	return fmt.Sprintf("%s/e/%s/%s/%s", schoolId, version, "", fMediaPath), fmt.Sprintf("%s/e/%s/%s/%s", schoolId, version, "", fMetaPath)
}
