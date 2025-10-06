package usecase

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AttachmentUseCase struct {
	S3     *s3.Client
	Bucket string
}

func (fs *AttachmentUseCase) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	key := header.Filename

	_, err = fs.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &fs.Bucket,
		Key:    &key,
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		fmt.Printf("failed to upload file to S3: %v", err)
		return "", err
	}

	return key, nil
}

func (fs *AttachmentUseCase) DownloadFile(key string) ([]byte, error) {
	resp, err := fs.S3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &fs.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fs *AttachmentUseCase) DeleteFile(key string) error {
	_, err := fs.S3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &fs.Bucket,
		Key:    &key,
	})
	return err
}
